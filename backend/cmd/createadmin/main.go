// createadmin: ilk yönetici kullanıcısını oluşturmak için tek seferlik CLI aracı.
// Kullanım: go run ./cmd/createadmin -username=admin -password=guclu-bir-sifre -role=admin
package main

import (
	"flag"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mezun-anket-backend/internal/config"
	"mezun-anket-backend/internal/domain"
)

func main() {
	username := flag.String("username", "", "yönetici kullanıcı adı")
	password := flag.String("password", "", "yönetici şifresi")
	role := flag.String("role", "admin", "admin | viewer")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("kullanım: go run ./cmd/createadmin -username=... -password=... -role=admin")
	}

	cfg := config.Load()
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&domain.AdminUser{}); err != nil {
		log.Fatal(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	admin := domain.AdminUser{Username: *username, PasswordHash: string(hash), Role: *role}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatal("kullanıcı oluşturulamadı: ", err)
	}
	log.Printf("yönetici kullanıcı oluşturuldu: %s (rol: %s)\n", *username, *role)
}
