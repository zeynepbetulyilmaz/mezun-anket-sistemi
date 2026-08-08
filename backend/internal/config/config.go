package config

import (
	"os"
)

type Config struct {
	Port           string
	DatabaseDSN    string
	JWTSecret      string
	EncryptionKey  string
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	SMTPFrom       string
	AllowedOrigins string
	Env            string
	InviteBaseURL  string
}

func Load() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseDSN:    getEnv("DATABASE_DSN", "host=postgres-db user=postgres password=password dbname=mezun_anket_db port=5432 sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "CHANGE_ME_IN_PRODUCTION"),
		EncryptionKey:  getEnv("ENCRYPTION_KEY", ""),
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