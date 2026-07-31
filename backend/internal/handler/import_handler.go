package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mezun-anket-backend/internal/apperror"
	"mezun-anket-backend/internal/middleware"
	"mezun-anket-backend/internal/service"
)

type ImportHandler struct {
	importSvc *service.ImportService
	mailSvc   *service.MailService
	inviteURL string // ör. https://anket.mersin.edu.tr/giris
}

func NewImportHandler(i *service.ImportService, m *service.MailService, inviteURL string) *ImportHandler {
	return &ImportHandler{importSvc: i, mailSvc: m, inviteURL: inviteURL}
}

type importRequest struct {
	Graduates   []service.GraduateImportRow `json:"graduates" binding:"required"`
	SendInvites bool                        `json:"sendInvites"`
}

// POST /api/v1/admin/graduates/import
// OBS'den gelen mezun verisini (soyad hariç, id hash string olarak) toplu içe aktarır.
func (h *ImportHandler) Import(c *gin.Context) {
	var req importRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.Validation("Geçersiz içe aktarma verisi."))
		return
	}

	result, tokens, err := h.importSvc.ImportBatch(req.Graduates, h.inviteURL)
	if err != nil {
		middleware.Fail(c, apperror.Internal("İçe aktarma sırasında hata oluştu."))
		return
	}

	if req.SendInvites {
		for _, row := range req.Graduates {
			if link, ok := tokens[row.OBSHashID]; ok {
				_ = h.mailSvc.EnqueueInviteEmail(row.Email, row.FirstName, link)
			}
		}
	}

	c.JSON(http.StatusOK, middleware.SuccessEnvelope(gin.H{
		"result": result,
		// inviteLinks: sadece bu response'ta bir kere döner (ham token DB'de tutulmaz,
		// yalnızca hash'i tutulur) - admin isterse manuel de paylaşabilsin diye.
		"inviteLinks": tokens,
	}))
}
