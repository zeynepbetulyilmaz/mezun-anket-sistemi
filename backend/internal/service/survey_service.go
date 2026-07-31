package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mezun-anket-backend/internal/apperror"
	"mezun-anket-backend/internal/crypto"
	"mezun-anket-backend/internal/domain"
)

type SurveyService struct {
	db  *gorm.DB
	enc *crypto.Encryptor
}

func NewSurveyService(db *gorm.DB, enc *crypto.Encryptor) *SurveyService {
	return &SurveyService{db: db, enc: enc}
}

// GetStructure: 5 kategori + soruları, adminin girdiği sıraya göre döner.
// Sorular DB'de tutulur ki ileride kod değiştirmeden güncellenebilsin.
func (s *SurveyService) GetStructure() ([]domain.SurveyCategory, error) {
	var categories []domain.SurveyCategory
	err := s.db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("survey_questions.order ASC")
	}).Order("\"order\" ASC").Find(&categories).Error
	return categories, err
}

// GetOrCreateResponse: mezunun anket kaydı yoksa oluşturur (idempotent),
// varsa mevcut cevaplarla birlikte döner (kaldığı yerden devam edebilsin).
func (s *SurveyService) GetOrCreateResponse(graduateID uint) (*domain.SurveyResponse, []domain.SurveyAnswer, error) {
	var resp domain.SurveyResponse
	err := s.db.Where("graduate_id = ?", graduateID).First(&resp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp = domain.SurveyResponse{
			GraduateID:  graduateID,
			Status:      "in_progress",
			CurrentStep: 1,
			StartedAt:   time.Now(),
		}
		if err := s.db.Create(&resp).Error; err != nil {
			return nil, nil, err
		}
	} else if err != nil {
		return nil, nil, err
	}

	var answers []domain.SurveyAnswer
	if err := s.db.Where("response_id = ?", resp.ID).Find(&answers).Error; err != nil {
		return nil, nil, err
	}
	return &resp, answers, nil
}

type StepAnswerInput struct {
	QuestionID uint   `json:"questionId" binding:"required"`
	Value      string `json:"value"`
}

// SaveStep: bir adımdaki tüm cevapları upsert eder ve CurrentStep'i günceller.
// Autosave mantığı: her "İleri" tıklamasında çağrılır, kullanıcı yarıda
// bırakırsa bir dahaki girişte kaldığı adımdan devam eder.
func (s *SurveyService) SaveStep(graduateID uint, stepNo int, answers []StepAnswerInput, requiredQuestionIDs map[uint]bool) error {
	if len(answers) == 0 {
		return apperror.Validation("Bu adımda en az bir cevap bekleniyor.")
	}

	var resp domain.SurveyResponse
	if err := s.db.Where("graduate_id = ?", graduateID).First(&resp).Error; err != nil {
		return apperror.NotFound("Anket kaydı bulunamadı.")
	}
	if resp.Status == "completed" {
		return apperror.Conflict("Anket zaten tamamlanmış.")
	}

	// Zorunlu soru kontrolü
	answered := map[uint]bool{}
	var details []apperror.Detail
	for _, a := range answers {
		if a.Value != "" {
			answered[a.QuestionID] = true
		}
	}
	for qID, required := range requiredQuestionIDs {
		if required && !answered[qID] {
			details = append(details, apperror.Detail{Field: "questionId:" + itoa(qID), Message: "Bu soru zorunludur."})
		}
	}
	if len(details) > 0 {
		return apperror.Validation("Bazı zorunlu sorular boş bırakılmış.", details...)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, a := range answers {
			row := domain.SurveyAnswer{
				ResponseID: resp.ID,
				QuestionID: a.QuestionID,
				ValueText:  a.Value,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "response_id"}, {Name: "question_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"value_text", "updated_at"}),
			}).Create(&row).Error; err != nil {
				return err
			}
		}
		nextStep := stepNo + 1
		if nextStep > resp.CurrentStep {
			return tx.Model(&resp).Update("current_step", nextStep).Error
		}
		return nil
	})
}

// Complete: son adım onaylandığında anketi tamamlanmış işaretler ve
// teşekkür mailini DB kuyruğuna (email_outbox) ekler - gerçek SMTP
// gönderimi ayrı bir worker goroutine tarafından asenkron yapılır.
func (s *SurveyService) Complete(graduateID uint, mailSvc *MailService) error {
	var resp domain.SurveyResponse
	if err := s.db.Where("graduate_id = ?", graduateID).First(&resp).Error; err != nil {
		return apperror.NotFound("Anket kaydı bulunamadı.")
	}
	if resp.Status == "completed" {
		return apperror.Conflict("Anket zaten tamamlanmış.")
	}

	var grad domain.Graduate
	if err := s.db.First(&grad, graduateID).Error; err != nil {
		return apperror.NotFound("Mezun kaydı bulunamadı.")
	}

	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&resp).Updates(map[string]any{
			"status":       "completed",
			"completed_at": now,
		}).Error; err != nil {
			return err
		}
		return mailSvc.EnqueueThankYouEmail(tx, grad)
	})
}

func itoa(v uint) string {
	if v == 0 {
		return "0"
	}
	digits := ""
	for v > 0 {
		digits = string(rune('0'+v%10)) + digits
		v /= 10
	}
	return digits
}
