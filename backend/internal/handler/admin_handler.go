package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mezun-anket-backend/internal/apperror"
	"mezun-anket-backend/internal/middleware"
	"mezun-anket-backend/internal/service"
)

type AdminHandler struct {
	admin *service.AdminService
}

func NewAdminHandler(a *service.AdminService) *AdminHandler {
	return &AdminHandler{admin: a}
}

// GET /api/v1/admin/stats/overview
func (h *AdminHandler) Overview(c *gin.Context) {
	stats, err := h.admin.Overview()
	if err != nil {
		middleware.Fail(c, apperror.Internal("İstatistikler yüklenemedi."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(stats))
}

// GET /api/v1/admin/stats/question/:code
// Örn: /api/v1/admin/stats/question/Q21 -> sektör dağılımı (pasta grafik verisi)
func (h *AdminHandler) QuestionDistribution(c *gin.Context) {
	code := c.Param("code")
	items, err := h.admin.QuestionDistribution(code)
	if err != nil {
		middleware.Fail(c, apperror.Internal("Soru dağılımı yüklenemedi."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(items))
}
