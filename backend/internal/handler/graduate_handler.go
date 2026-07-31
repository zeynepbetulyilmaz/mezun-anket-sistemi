package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mezun-anket-backend/internal/apperror"
	"mezun-anket-backend/internal/middleware"
	"mezun-anket-backend/internal/service"
)

type GraduateHandler struct {
	graduates *service.GraduateService
}

func NewGraduateHandler(g *service.GraduateService) *GraduateHandler {
	return &GraduateHandler{graduates: g}
}

// GET /api/v1/me
func (h *GraduateHandler) Me(c *gin.Context) {
	graduateID := c.MustGet(middleware.CtxGraduateID).(uint)
	grad, err := h.graduates.Me(graduateID)
	if err != nil {
		middleware.Fail(c, apperror.NotFound("Mezun kaydı bulunamadı."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(gin.H{
		"firstName":      grad.FirstName,
		"facultyName":    grad.FacultyName,
		"departmentName": grad.DepartmentName,
		"graduationYear": grad.GraduationYear,
	}))
}
