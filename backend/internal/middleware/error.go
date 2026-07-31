package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mezun-anket-backend/internal/apperror"
)

// SuccessEnvelope: tüm başarılı response'ları saran standart zarf.
func SuccessEnvelope(data interface{}) gin.H {
	return gin.H{"success": true, "data": data}
}

// ErrorEnvelope: gin.Context.Error(...) ile eklenen hataları yakalayıp
// standart {success:false, error:{...}} formatına çevirir.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if appErr, ok := err.(*apperror.AppError); ok {
			c.JSON(appErr.HTTPStatus, gin.H{"success": false, "error": appErr})
			return
		}

		// Beklenmeyen hata: detay dışa sızdırılmaz, sadece loglanır.
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Beklenmeyen bir hata oluştu.",
			},
		})
	}
}

// Fail: handler'larda kısa yoldan hata döndürmek için yardımcı.
func Fail(c *gin.Context, err *apperror.AppError) {
	_ = c.Error(err)
	c.Abort()
}
