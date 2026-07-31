package config

import (
	"os"
)

// Config uygulamanın ortam değişkenlerinden okunan tüm ayarlarını tutar.
// Hiçbir gizli değer koda gömülmez; production'da bunlar secret manager
// (ör. Vault, Doppler, K8s Secret) üzerinden env'e enjekte edilmelidir.
type Config struct {
	Port           string
	DatabaseDSN    string
	JWTSecret      string
	EncryptionKey  string // 32 byte (AES-256) - hex ya da base64 olarak env'de tutulup decode edilir
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	SMTPFrom       string
	AllowedOrigins string
	Env            string // "development" | "production"
	InviteBaseURL  string // ör. https://anket.mersin.edu.tr/giris
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseDSN:    getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=mezun_anket port=5432 sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "CHANGE_ME_IN_PRODUCTION"),
		EncryptionKey:  getEnv("ENCRYPTION_KEY", ""), // boşsa main.go başlangıçta panik atar
		SMTPHost:       getEnv("SMTP_HOST", "smtp.mersin.edu.tr"),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:       getEnv("SMTP_FROM", "mezunanket@mersin.edu.tr"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173"),
		Env:            getEnv("APP_ENV", "development"),
		InviteBaseURL:  getEnv("INVITE_BASE_URL", "http://localhost:5173/giris"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
