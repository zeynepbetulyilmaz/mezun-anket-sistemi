package service

import (
	"gorm.io/gorm"

	"mezun-anket-backend/internal/domain"
)

type GraduateService struct {
	db *gorm.DB
}

func NewGraduateService(db *gorm.DB) *GraduateService {
	return &GraduateService{db: db}
}

// Me: kişiselleştirilmiş karşılama ekranı için ad/bölüm/mezuniyet yılı döner.
// İletişim bilgileri (şifreli tutulan e-posta/telefon) bu endpoint'te ASLA
// decrypt edilip dönülmez.
func (g *GraduateService) Me(graduateID uint) (*domain.Graduate, error) {
	var grad domain.Graduate
	if err := g.db.Select("id", "first_name", "faculty_name", "department_name", "graduation_year").
		First(&grad, graduateID).Error; err != nil {
		return nil, err
	}
	return &grad, nil
}
