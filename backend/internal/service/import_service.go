package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"

	"mezun-anket-backend/internal/crypto"
	"mezun-anket-backend/internal/domain"
)

type ImportService struct {
	db  *gorm.DB
	enc *crypto.Encryptor
}

func NewImportService(db *gorm.DB, enc *crypto.Encryptor) *ImportService {
	return &ImportService{db: db, enc: enc}
}

// GraduateImportRow: OBS'nin gönderdiği/öngörülen format.
// Önemli: "soyad" alanı kasıtlı olarak burada yok — OBS zaten göndermiyor.
type GraduateImportRow struct {
	OBSHashID      string `json:"obsHashId" binding:"required"`
	FirstName      string `json:"firstName" binding:"required"`
	FacultyName    string `json:"facultyName"`
	DepartmentName string `json:"departmentName"`
	GraduationYear int    `json:"graduationYear"`
	StudentNoHash  string `json:"studentNoHash"`
	Email          string `json:"email"` // düz metin gelir, DB'ye yazmadan önce şifrelenir
	Phone          string `json:"phone"`
}

type ImportResult struct {
	Inserted int      `json:"inserted"`
	Updated  int      `json:"updated"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

// randomToken: mezuna gönderilecek tek kullanımlık giriş linki için rastgele token üretir.
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashTokenValue(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// ImportBatch: OBS'den gelen mezun listesini upsert eder, her biri için
// tek kullanımlık giriş token'ı üretir ve rawToken'ı (yalnızca bir kez,
// DB'ye hash'lenmiş hali yazılarak) çağırana döner ki davet e-postasına
// gömülsün. Aynı akışta davet e-postası da email_outbox'a (DB kuyruğu)
// eklenir.
func (s *ImportService) ImportBatch(rows []GraduateImportRow, inviteLinkBase string) (*ImportResult, map[string]string, error) {
	result := &ImportResult{}
	rawTokensByHash := map[string]string{} // OBSHashID -> rawToken (sadece bu response'ta döner)

	for _, row := range rows {
		var existing domain.Graduate
		err := s.db.Where("obs_hash_id = ?", row.OBSHashID).First(&existing).Error

		emailEnc, emailNonce, encErr := s.enc.Encrypt(row.Email)
		if encErr != nil {
			result.Failed++
			result.Errors = append(result.Errors, row.OBSHashID+": e-posta şifrelenemedi")
			continue
		}
		phoneEnc, phoneNonce, encErr2 := s.enc.Encrypt(row.Phone)
		if encErr2 != nil {
			result.Failed++
			result.Errors = append(result.Errors, row.OBSHashID+": telefon şifrelenemedi")
			continue
		}

		if err == gorm.ErrRecordNotFound {
			grad := domain.Graduate{
				OBSHashID:      row.OBSHashID,
				FirstName:      row.FirstName,
				FacultyName:    row.FacultyName,
				DepartmentName: row.DepartmentName,
				GraduationYear: row.GraduationYear,
				StudentNoHash:  row.StudentNoHash,
				EmailEnc:       emailEnc,
				EmailNonce:     emailNonce,
				PhoneEnc:       phoneEnc,
				PhoneNonce:     phoneNonce,
			}
			if err := s.db.Create(&grad).Error; err != nil {
				result.Failed++
				result.Errors = append(result.Errors, row.OBSHashID+": "+err.Error())
				continue
			}
			result.Inserted++
			existing = grad
		} else if err == nil {
			s.db.Model(&existing).Updates(map[string]any{
				"first_name":      row.FirstName,
				"faculty_name":    row.FacultyName,
				"department_name": row.DepartmentName,
				"graduation_year": row.GraduationYear,
				"student_no_hash": row.StudentNoHash,
				"email_enc":       emailEnc,
				"email_nonce":     emailNonce,
				"phone_enc":       phoneEnc,
				"phone_nonce":     phoneNonce,
			})
			result.Updated++
		} else {
			result.Failed++
			result.Errors = append(result.Errors, row.OBSHashID+": "+err.Error())
			continue
		}

		rawToken, tErr := randomToken()
		if tErr != nil {
			continue
		}
		loginToken := domain.LoginToken{
			GraduateID: existing.ID,
			TokenHash:  hashTokenValue(rawToken),
			ExpiresAt:  time.Now().Add(90 * 24 * time.Hour),
			CreatedAt:  time.Now(),
		}
		if err := s.db.Create(&loginToken).Error; err == nil {
			rawTokensByHash[row.OBSHashID] = inviteLinkBase + "?token=" + rawToken
		}
	}

	return result, rawTokensByHash, nil
}
