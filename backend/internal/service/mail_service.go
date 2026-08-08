package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"github.com/google/uuid"

	"mezun-anket-backend/internal/crypto"
	"mezun-anket-backend/internal/domain"
)

type MailService struct {
	db  *gorm.DB
	enc *crypto.Encryptor
}

func NewMailService(db *gorm.DB, enc *crypto.Encryptor) *MailService {
	return &MailService{db: db, enc: enc}
}

func (m *MailService) EnqueueThankYouEmail(tx *gorm.DB, grad domain.Graduate) error {
	plainEmail, err := m.enc.Decrypt(grad.EmailEnc, grad.EmailNonce)
	if err != nil || plainEmail == "" {
		return nil
	}
	cipherText, nonce, err := m.enc.Encrypt(plainEmail)
	if err != nil {
		return err
	}

	subject := "Anketiniz için teşekkür ederiz - Mersin Üniversitesi"
	body := fmt.Sprintf(
		"Sayın %s,\n\nMezun Takip Anketi'ni tamamladığınız için teşekkür ederiz. "+
			"Katkılarınız üniversitemizin eğitim kalitesini geliştirmesine yardımcı olacaktır.\n\n"+
			"Mersin Üniversitesi Mezun İlişkileri",
		grad.FirstName,
	)

	outbox := domain.EmailOutbox{
		ToEmailEnc:   cipherText,
		ToEmailNonce: nonce,
		Subject:      subject,
		Body:         body,
		Status:       "pending",
		SendAfter:    time.Now(),
		CreatedAt:    time.Now(),
	}
	return tx.Create(&outbox).Error
}

func (m *MailService) EnqueueInviteEmail(plainEmail, firstName, inviteLink string) error {
	if plainEmail == "" {
		return nil
	}
	cipherText, nonce, err := m.enc.Encrypt(plainEmail)
	if err != nil {
		return err
	}
	subject := "Mezun Takip Anketimize Katılın - Mersin Üniversitesi"
	body := fmt.Sprintf(
		"Sayın %s,\n\nMersin Üniversitesi Mezun Takip Anketimize aşağıdaki kişisel linkten katılabilirsiniz:\n%s\n\nKatkılarınız için şimdiden teşekkür ederiz.\n\nMersin Üniversitesi Mezun İlişkileri",
		firstName, inviteLink,
	)
	outbox := domain.EmailOutbox{
		ToEmailEnc:   cipherText,
		ToEmailNonce: nonce,
		Subject:      subject,
		Body:         body,
		Status:       "pending",
		SendAfter:    time.Now(),
		CreatedAt:    time.Now(),
	}
	return m.db.Create(&outbox).Error
}

func (m *MailService) QueueInvites(req domain.SendInvitesRequest) (int, error) {
	var graduates []domain.Graduate
	query := m.db.Model(&domain.Graduate{})

	if !req.SendToAll {
		if req.DepartmentName == "" {
			return 0, fmt.Errorf("lütfen bir bölüm seçiniz veya 'Tümüne Gönder'i işaretleyiniz")
		}
		query = query.Where("department_name = ?", req.DepartmentName)
	}

	if err := query.Find(&graduates).Error; err != nil {
		return 0, err
	}

	queuedCount := 0
	err := m.db.Transaction(func(tx *gorm.DB) error {
		for _, grad := range graduates {
			plainEmail, err := m.enc.Decrypt(grad.EmailEnc, grad.EmailNonce)
			if err != nil || plainEmail == "" {
				continue
			}

			cipherText, nonce, err := m.enc.Encrypt(plainEmail)
			if err != nil {
				return err
			}

			// HATA BURADAYDI: Olmayan survey_token sütunu yerine, LoginToken tablosunu kullanıyoruz!
			token := uuid.New().String()
			loginToken := domain.LoginToken{
				GraduateID: grad.ID,
				TokenHash:  token, // Linkin eşsiz token'ı
				ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // Link 1 hafta geçerli olsun
			}
			
			if err := tx.Create(&loginToken).Error; err != nil {
				return err
			}

			surveyLink := fmt.Sprintf("http://localhost:5173/welcome?token=%s", token)

			subject := "Mersin Üniversitesi - Mezun Bilgi Anketi Daveti"
			body := fmt.Sprintf("Değerli mezunumuz,\n\nEğitim kalitemizi artırmak ve kariyer gelişiminizi takip etmek amacıyla hazırladığımız ankete katılmak için aşağıdaki linke tıklayınız:\n\n%s\n\nTeşekkürler.", surveyLink)

			outbox := domain.EmailOutbox{
				ToEmailEnc:   cipherText,
				ToEmailNonce: nonce,
				Subject:      subject,
				Body:         body,
				Status:       "pending",
				SendAfter:    time.Now(),
				CreatedAt:    time.Now(),
			}

			if err := tx.Create(&outbox).Error; err != nil {
				return err
			}
			queuedCount++
		}
		return nil
	})

	return queuedCount, err
}
