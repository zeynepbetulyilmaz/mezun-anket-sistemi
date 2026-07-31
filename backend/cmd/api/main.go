package main

import (
	"log"
    "github.com/joho/godotenv"
	"context"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mezun-anket-backend/internal/config"
	"mezun-anket-backend/internal/crypto"
	"mezun-anket-backend/internal/domain"
	"mezun-anket-backend/internal/handler"
	"mezun-anket-backend/internal/mail"
	"mezun-anket-backend/internal/middleware"
	"mezun-anket-backend/internal/seed"
	"mezun-anket-backend/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Println("Uyarı: .env dosyası bulunamadı, sistem değişkenleri kullanılacak.")
    }
	cfg := config.Load()

	if cfg.EncryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY tanımlı değil. `openssl rand -hex 32` ile üretip ortam değişkeni olarak ayarlayın.")
	}
	encryptor, err := crypto.NewEncryptor(cfg.EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("veritabanına bağlanılamadı: %v", err)
	}

	if err := db.AutoMigrate(domain.AllModels()...); err != nil {
		log.Fatalf("migration hatası: %v", err)
	}

	if err := seed.Run(db); err != nil {
		log.Fatalf("seed hatası: %v", err)
	}

	// --- Servisler ---
	authSvc := service.NewAuthService(db, cfg.JWTSecret)
	graduateSvc := service.NewGraduateService(db)
	surveySvc := service.NewSurveyService(db, encryptor)
	mailSvc := service.NewMailService(db, encryptor)
	adminSvc := service.NewAdminService(db)
	importSvc := service.NewImportService(db, encryptor)

	// --- Handler'lar ---
	authHandler := handler.NewAuthHandler(authSvc)
	graduateHandler := handler.NewGraduateHandler(graduateSvc)
	surveyHandler := handler.NewSurveyHandler(surveySvc, mailSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)
	importHandler := handler.NewImportHandler(importSvc, mailSvc, cfg.InviteBaseURL)

	// --- Mail worker: harici kuyruk teknolojisi olmadan, DB üzerinde ---
	smtpClient := mail.NewSMTPClient(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFrom)
	worker := mail.NewWorker(db, smtpClient, encryptor)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go worker.Run(ctx)

	// --- Router ---
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(cfg.AllowedOrigins))
	r.Use(middleware.ErrorHandler())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	{
		api.POST("/auth/token-login", authHandler.TokenLogin)
		api.POST("/admin/login", authHandler.AdminLogin)

		// Mezun tarafı - JWT zorunlu
		grad := api.Group("")
		grad.Use(middleware.RequireGraduateAuth(cfg.JWTSecret))
		{
			grad.GET("/me", graduateHandler.Me)
			grad.GET("/survey/structure", surveyHandler.Structure)
			grad.GET("/survey/response", surveyHandler.GetResponse)
			grad.PUT("/survey/response/step/:stepNo", surveyHandler.SaveStep)
			grad.POST("/survey/response/complete", surveyHandler.Complete)
		}

		// Admin tarafı - JWT + rol zorunlu
		admin := api.Group("/admin")
		admin.Use(middleware.RequireAdminAuth(cfg.JWTSecret, "admin", "viewer"))
		{
			admin.GET("/stats/overview", adminHandler.Overview)
			admin.GET("/stats/question/:code", adminHandler.QuestionDistribution)
			admin.POST("/graduates/import", importHandler.Import)
		}
	}

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	go func() {
		log.Println("[api] dinleniyor :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("sunucu hatası: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("[api] kapatılıyor...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

// corsMiddleware: ALLOWED_ORIGINS env'i virgülle ayrılmış birden fazla origin destekler.
func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")
	allowed := map[string]bool{}
	for _, o := range origins {
		allowed[strings.TrimSpace(o)] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
