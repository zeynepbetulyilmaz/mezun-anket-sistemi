package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mezun-anket-backend/internal/apperror"
	"mezun-anket-backend/internal/middleware"
	"mezun-anket-backend/internal/service"
)

type SurveyHandler struct {
	survey *service.SurveyService
	mail   *service.MailService
}

func NewSurveyHandler(survey *service.SurveyService, mail *service.MailService) *SurveyHandler {
	return &SurveyHandler{survey: survey, mail: mail}
}

// GET /api/v1/survey/structure
func (h *SurveyHandler) Structure(c *gin.Context) {
	categories, err := h.survey.GetStructure()
	if err != nil {
		middleware.Fail(c, apperror.Internal("Anket yapısı yüklenemedi."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(categories))
}

// GET /api/v1/survey/response
func (h *SurveyHandler) GetResponse(c *gin.Context) {
	graduateID := c.MustGet(middleware.CtxGraduateID).(uint)
	resp, answers, err := h.survey.GetOrCreateResponse(graduateID)
	if err != nil {
		middleware.Fail(c, apperror.Internal("Anket kaydı alınamadı."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(gin.H{
		"response": resp,
		"answers":  answers,
	}))
}

type saveStepRequest struct {
	Answers []service.StepAnswerInput `json:"answers" binding:"required"`
	// RequiredQuestionIds: frontend'in o adımda zorunlu gösterdiği soru id'leri.
	// (Ana doğrulama kaynağı DB'deki SurveyQuestion.Required olsa da,
	// handler basit tutulması için burada da kabul eder.)
	RequiredQuestionIds []uint `json:"requiredQuestionIds"`
}

// PUT /api/v1/survey/response/step/:stepNo
func (h *SurveyHandler) SaveStep(c *gin.Context) {
	graduateID := c.MustGet(middleware.CtxGraduateID).(uint)

	stepNo, err := strconv.Atoi(c.Param("stepNo"))
	if err != nil || stepNo < 1 || stepNo > 6 {
		middleware.Fail(c, apperror.Validation("Geçersiz adım numarası (1-5 olmalı)."))
		return
	}

	var req saveStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.Validation("İstek gövdesi hatalı."))
		return
	}

	required := map[uint]bool{}
	for _, id := range req.RequiredQuestionIds {
		required[id] = true
	}

	if err := h.survey.SaveStep(graduateID, stepNo, req.Answers, required); err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			middleware.Fail(c, appErr)
			return
		}
		middleware.Fail(c, apperror.Internal("Cevaplar kaydedilemedi."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(gin.H{"saved": true}))
}

// POST /api/v1/survey/response/complete
func (h *SurveyHandler) Complete(c *gin.Context) {
	graduateID := c.MustGet(middleware.CtxGraduateID).(uint)
	if err := h.survey.Complete(graduateID, h.mail); err != nil {
		if appErr, ok := err.(*apperror.AppError); ok {
			middleware.Fail(c, appErr)
			return
		}
		middleware.Fail(c, apperror.Internal("Anket tamamlanamadı."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(gin.H{"completed": true}))
}
