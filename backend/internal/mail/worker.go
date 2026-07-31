// Package mail: harici kuyruk teknolojisi (Redis/RabbitMQ/BullMQ vb.) kullanmadan,
// PostgreSQL'in "SELECT ... FOR UPDATE SKIP LOCKED" özelliğiyle güvenli, birden
// fazla instance'da bile çakışmayan asenkron mail işleme sağlar.
package mail

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"mezun-anket-backend/internal/crypto"
	"mezun-anket-backend/internal/domain"
)

type Worker struct {
	db       *gorm.DB
	smtp     *SMTPClient
	enc      *crypto.Encryptor
	interval time.Duration
	batch    int
}

func NewWorker(db *gorm.DB, smtpClient *SMTPClient, enc *crypto.Encryptor) *Worker {
	return &Worker{db: db, smtp: smtpClient, enc: enc, interval: 15 * time.Second, batch: 20}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	log.Println("[mail-worker] başlatıldı, her", w.interval, "saniyede bir email_outbox taranacak")

	for {
		select {
		case <-ctx.Done():
			log.Println("[mail-worker] durduruldu")
			return
		case <-ticker.C:
			if err := w.processBatch(); err != nil {
				log.Println("[mail-worker] batch hatası:", err)
			}
		}
	}
}

func (w *Worker) processBatch() error {
	var jobs []domain.EmailOutbox

	err := w.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw(`
			SELECT * FROM email_outboxes
			WHERE status = 'pending' AND send_after <= now()
			ORDER BY id
			LIMIT ?
			FOR UPDATE SKIP LOCKED
		`, w.batch).Scan(&jobs).Error; err != nil {
			return err
		}
		if len(jobs) == 0 {
			return nil
		}
		ids := make([]uint, len(jobs))
		for i, j := range jobs {
			ids[i] = j.ID
		}
		return tx.Model(&domain.EmailOutbox{}).Where("id IN ?", ids).
			Updates(map[string]any{"status": "processing", "locked_at": time.Now()}).Error
	})
	if err != nil {
		return err
	}

	for _, job := range jobs {
		email, decErr := w.enc.Decrypt(job.ToEmailEnc, job.ToEmailNonce)
		if decErr != nil || email == "" {
			w.markFailed(job, "adres çözülemedi")
			continue
		}
		if sendErr := w.smtp.Send(email, job.Subject, job.Body); sendErr != nil {
			w.retry(job, sendErr.Error())
			continue
		}
		now := time.Now()
		w.db.Model(&domain.EmailOutbox{}).Where("id = ?", job.ID).
			Updates(map[string]any{"status": "sent", "sent_at": now})
	}
	return nil
}

func (w *Worker) retry(job domain.EmailOutbox, errMsg string) {
	attempts := job.Attempts + 1
	backoff := time.Duration(attempts*attempts) * time.Minute // basit exponential backoff
	if attempts >= 5 {
		w.markFailed(job, errMsg)
		return
	}
	w.db.Model(&domain.EmailOutbox{}).Where("id = ?", job.ID).Updates(map[string]any{
		"status":     "pending",
		"attempts":   attempts,
		"last_error": errMsg,
		"send_after": time.Now().Add(backoff),
	})
}

func (w *Worker) markFailed(job domain.EmailOutbox, errMsg string) {
	w.db.Model(&domain.EmailOutbox{}).Where("id = ?", job.ID).Updates(map[string]any{
		"status":     "failed",
		"last_error": errMsg,
	})
}
