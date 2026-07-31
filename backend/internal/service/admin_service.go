package service

import (
	"gorm.io/gorm"
)

type AdminService struct {
	db *gorm.DB
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{db: db}
}

type OverviewStats struct {
	TotalGraduates    int64            `json:"totalGraduates"`
	TotalResponses    int64            `json:"totalResponses"`
	CompletedCount    int64            `json:"completedCount"`
	CompletionRate    float64          `json:"completionRate"`
	DropOffByStep     map[int]int64    `json:"dropOffByStep"` // hangi adımda kaç kişi yarıda bıraktı
}

// Overview: yalnızca sayısal/istatistiksel veri döner - hiçbir kişisel veri
// (ad, e-posta, telefon, ham cevap metni) frontend'e sızmaz.
func (a *AdminService) Overview() (*OverviewStats, error) {
	var totalGraduates, totalResponses, completed int64
	a.db.Table("graduates").Count(&totalGraduates)
	a.db.Table("survey_responses").Count(&totalResponses)
	a.db.Table("survey_responses").Where("status = ?", "completed").Count(&completed)

	rate := 0.0
	if totalResponses > 0 {
		rate = float64(completed) / float64(totalResponses) * 100
	}

	type stepCount struct {
		CurrentStep int
		Cnt         int64
	}
	var rows []stepCount
	a.db.Table("survey_responses").
		Select("current_step, count(*) as cnt").
		Where("status = ?", "in_progress").
		Group("current_step").
		Scan(&rows)

	dropOff := map[int]int64{}
	for _, r := range rows {
		dropOff[r.CurrentStep] = r.Cnt
	}

	return &OverviewStats{
		TotalGraduates: totalGraduates,
		TotalResponses: totalResponses,
		CompletedCount: completed,
		CompletionRate: rate,
		DropOffByStep:  dropOff,
	}, nil
}

type DistributionItem struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

// QuestionDistribution: herhangi bir sorunun (kod ile, ör. "Q21") cevap
// dağılımını döner - Recharts Pie/Bar için doğrudan tüketilebilir format.
func (a *AdminService) QuestionDistribution(questionCode string) ([]DistributionItem, error) {
	var items []DistributionItem
	err := a.db.Table("survey_answers a").
		Joins("JOIN survey_questions q ON q.id = a.question_id").
		Select("a.value_text as label, count(*) as count").
		Where("q.code = ?", questionCode).
		Group("a.value_text").
		Order("count DESC").
		Scan(&items).Error
	return items, err
}
