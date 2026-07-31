package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

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

// EnqueueThankYouEmail: gerçek SMTP gönderimini YAPMAZ; sadece email_outbox
// tablosuna (DB-üzeri kuyruk) bir kayıt düşer. Gönderim MailWorker tarafından
// asenkron olarak, ayrı bir goroutine'de yapılır.
func (m *MailService) EnqueueThankYouEmail(tx *gorm.DB, grad domain.Graduate) error {
	plainEmail, err := m.enc.Decrypt(grad.EmailEnc, grad.EmailNonce)
	if err != nil || plainEmail == "" {
		// E-posta çözülemiyorsa anketin tamamlanmasını engellemeyelim,
		// sadece mail kuyruklanmaz; bu durum loglanmalı.
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

// EnqueueInviteEmail: OBS import sonrası mezuna anket davet linkini içeren
// e-postayı DB kuyruğuna (email_outbox) ekler.
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
