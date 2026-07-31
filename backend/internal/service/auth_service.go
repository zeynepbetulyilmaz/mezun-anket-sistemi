package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mezun-anket-backend/internal/domain"
	"mezun-anket-backend/internal/middleware"
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{db: db, jwtSecret: jwtSecret}
}

// hashToken: ham token DB'de tutulmaz, SHA-256 hash'i tutulur (tek kullanımlık
// linki başkası DB dump'ından okuyup tekrar kullanamasın diye).
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

var (
	ErrTokenInvalid = errors.New("token geçersiz veya süresi dolmuş")
	ErrTokenUsed    = errors.New("bu link daha önce kullanılmış")
)

// LoginWithToken: OBS/mail ile mezuna gönderilen tek kullanımlık linkteki
// token'ı doğrular, JWT üretir. Not: LoginToken kaydı "tek kullanımlıktır"
// ama link tekrar açıldığında mezun anketine kaldığı yerden devam edebilsin
// diye kullanım sonrası JWT'nin süresi yeterince uzun tutulur (ör. 30 gün),
// tek kullanımlık kısıtı sadece "linkin kendisinin" tekrar tekrar
// paylaşılmasını engellemek içindir.
func (s *AuthService) LoginWithToken(rawToken string) (jwtStr string, graduate *domain.Graduate, err error) {
	hashed := hashToken(rawToken)

	var lt domain.LoginToken
	if err := s.db.Where("token_hash = ?", hashed).First(&lt).Error; err != nil {
		return "", nil, ErrTokenInvalid
	}
	if lt.ExpiresAt.Before(time.Now()) {
		return "", nil, ErrTokenInvalid
	}

	var grad domain.Graduate
	if err := s.db.First(&grad, lt.GraduateID).Error; err != nil {
		return "", nil, ErrTokenInvalid
	}

	if lt.UsedAt == nil {
		now := time.Now()
		s.db.Model(&lt).Update("used_at", now)
	}

	claims := jwt.MapClaims{
		middleware.ClaimGraduateID: grad.ID,
		"exp":                      time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":                      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return signed, &grad, nil
}

// AdminLogin: kullanıcı adı + bcrypt şifre doğrulaması, JWT üretir.
func (s *AuthService) AdminLogin(username, password string) (string, *domain.AdminUser, error) {
	var admin domain.AdminUser
	if err := s.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return "", nil, errors.New("kullanıcı adı veya şifre hatalı")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("kullanıcı adı veya şifre hatalı")
	}
	claims := jwt.MapClaims{
		middleware.ClaimAdminID: admin.ID,
		middleware.ClaimRole:    admin.Role,
		"exp":                   time.Now().Add(8 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return signed, &admin, nil
}
