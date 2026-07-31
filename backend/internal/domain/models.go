package domain

import "time"

// Graduate: OBS'den senkronize edilen mezun kaydı.
// ÖNEMLİ: soyad OBS'den gelmiyor ve bu şemada hiç yer almıyor (veri minimizasyonu).
type Graduate struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	OBSHashID      string `gorm:"uniqueIndex;size:128;not null" json:"-"`
	FirstName      string `gorm:"size:100;not null" json:"firstName"`
	FacultyName    string `gorm:"size:150" json:"facultyName"`
	DepartmentName string `gorm:"size:150" json:"departmentName"`
	GraduationYear int    `json:"graduationYear"`
	StudentNoHash  string `gorm:"size:128;index" json:"-"`

	EmailEnc   []byte `gorm:"type:bytea" json:"-"`
	EmailNonce []byte `gorm:"type:bytea" json:"-"`
	PhoneEnc   []byte `gorm:"type:bytea" json:"-"`
	PhoneNonce []byte `gorm:"type:bytea" json:"-"`

	ConsentGivenAt *time.Time `json:"-"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// LoginToken: mezuna gönderilen tek kullanımlık giriş linkinin karşılığı.
// JWT'nin kendisi imzalı ve kısa ömürlü olsa da, "tek kullanımlık" garantisi
// için DB'de ayrıca izleniyor (opsiyonel ama önerilir).
type LoginToken struct {
	ID         uint   `gorm:"primaryKey"`
	GraduateID uint   `gorm:"index;not null"`
	TokenHash  string `gorm:"uniqueIndex;size:128;not null"`
	ExpiresAt  time.Time
	UsedAt     *time.Time
	CreatedAt  time.Time
}

type SurveyCategory struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Order int    `gorm:"not null" json:"order"`
	Title string `gorm:"size:150" json:"title"`

	Questions []SurveyQuestion `gorm:"foreignKey:CategoryID" json:"questions,omitempty"`
}

// AnswerType değerleri: scale_1_10 | single_choice | multi_choice | text | number | duration_months
type SurveyQuestion struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	CategoryID  uint   `gorm:"index;not null" json:"categoryId"`
	Order       int    `gorm:"not null" json:"order"`
	Code        string `gorm:"size:10;uniqueIndex" json:"code"`
	Text        string `gorm:"type:text" json:"text"`
	AnswerType  string `gorm:"size:20" json:"answerType"`
	OptionsJSON string `gorm:"type:jsonb" json:"optionsJson,omitempty"`
	Required    bool   `gorm:"default:true" json:"required"`
	TargetFaculty       *string `json:"targetFaculty" db:"target_faculty"`
	TargetDepartment    *string `json:"targetDepartment" db:"target_department"`
	DependsOnQuestionID *int64  `json:"dependsOnQuestionId" db:"depends_on_question_id"`
	DependsOnAnswer     *string `json:"dependsOnAnswer" db:"depends_on_answer"`
}

type SurveyResponse struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	GraduateID  uint       `gorm:"uniqueIndex;not null" json:"graduateId"`
	Status      string     `gorm:"size:20;default:'in_progress'" json:"status"` // in_progress | completed
	CurrentStep int        `gorm:"default:1" json:"currentStep"`
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
}

type SurveyAnswer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ResponseID uint      `gorm:"uniqueIndex:idx_resp_q;not null" json:"responseId"`
	QuestionID uint      `gorm:"uniqueIndex:idx_resp_q;not null" json:"questionId"`
	ValueText  string    `gorm:"type:text" json:"valueText"` // scale/number/text/choice hepsi burada normalize string/JSON olarak tutulur
	UpdatedAt  time.Time `json:"-"`
}

type AdminUser struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"uniqueIndex;size:100"`
	PasswordHash string `gorm:"size:255"`
	Role         string `gorm:"size:20"` // admin | viewer
	CreatedAt    time.Time
}

// EmailOutbox: harici kuyruk teknolojisi (Redis/RabbitMQ/BullMQ) kullanmadan
// PostgreSQL üzerinde tutulan asenkron mail kuyruğu.
type EmailOutbox struct {
	ID           uint       `gorm:"primaryKey"`
	ToEmailEnc   []byte     `gorm:"type:bytea"`
	ToEmailNonce []byte     `gorm:"type:bytea"`
	Subject      string     `gorm:"size:255"`
	Body         string     `gorm:"type:text"`
	Status       string     `gorm:"size:20;default:'pending';index"` // pending | processing | sent | failed
	Attempts     int        `gorm:"default:0"`
	LastError    string     `gorm:"type:text"`
	LockedAt     *time.Time
	SendAfter    time.Time `gorm:"index"`
	CreatedAt    time.Time
	SentAt       *time.Time
}

// AllModels: GORM AutoMigrate için tüm modellerin listesi.
func AllModels() []interface{} {
	return []interface{}{
		&Graduate{},
		&LoginToken{},
		&SurveyCategory{},
		&SurveyQuestion{},
		&SurveyResponse{},
		&SurveyAnswer{},
		&AdminUser{},
		&EmailOutbox{},
	}
}
