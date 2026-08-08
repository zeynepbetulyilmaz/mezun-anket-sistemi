package handler

import (
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	
	// Buradaki isimler "mezun-anket-backend" olarak düzeltildi:
	"mezun-anket-backend/internal/domain"
	"mezun-anket-backend/internal/service"
)

type AdminMailHandler struct {
	mailService *service.MailService
}

func NewAdminMailHandler(ms *service.MailService) *AdminMailHandler {
	return &AdminMailHandler{mailService: ms}
}

// SendInvites, React'ten gelen POST isteğini karşılar
func (h *AdminMailHandler) SendInvites(c *gin.Context) {
	var req domain.SendInvitesRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri formatı."})
		return
	}

	count, err := h.mailService.QueueInvites(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kuyruğa ekleme başarısız: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%d adet davetiye başarıyla mail kuyruğuna eklendi.", count),
		"queued":  count,
	})
}