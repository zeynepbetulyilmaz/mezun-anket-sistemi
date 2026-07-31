package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"mezun-anket-backend/internal/apperror"
)

const (
	ClaimGraduateID = "graduateId"
	ClaimAdminID    = "adminId"
	ClaimRole       = "role"

	CtxGraduateID = "ctx_graduate_id"
	CtxAdminID    = "ctx_admin_id"
	CtxRole       = "ctx_role"
)

// bearerToken: "Authorization: Bearer <token>" header'ından ham token'ı çıkarır.
func bearerToken(c *gin.Context) (string, bool) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(header, "Bearer "), true
}

// RequireGraduateAuth: mezun tarafı endpoint'lerini korur; JWT'den graduateId'yi
// context'e ekler.
func RequireGraduateAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, ok := bearerToken(c)
		if !ok {
			Fail(c, apperror.Unauthorized("Oturum bulunamadı."))
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			Fail(c, apperror.Unauthorized("Oturum süresi dolmuş, lütfen linki tekrar kullanın."))
			return
		}
		gradID, ok := claims[ClaimGraduateID].(float64)
		if !ok {
			Fail(c, apperror.Unauthorized("Geçersiz oturum bilgisi."))
			return
		}
		c.Set(CtxGraduateID, uint(gradID))
		c.Next()
	}
}

// RequireAdminAuth: admin panel endpoint'lerini korur, opsiyonel rol kontrolü yapar.
func RequireAdminAuth(secret string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, ok := bearerToken(c)
		if !ok {
			Fail(c, apperror.Unauthorized("Yönetici oturumu bulunamadı."))
			return
		}
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			Fail(c, apperror.Unauthorized("Oturum süresi dolmuş."))
			return
		}
		adminID, ok := claims[ClaimAdminID].(float64)
		if !ok {
			Fail(c, apperror.Unauthorized("Geçersiz oturum bilgisi."))
			return
		}
		role, _ := claims[ClaimRole].(string)
		if len(allowedRoles) > 0 {
			permitted := false
			for _, r := range allowedRoles {
				if r == role {
					permitted = true
					break
				}
			}
			if !permitted {
				Fail(c, apperror.Forbidden("Bu işlem için yetkiniz yok."))
				return
			}
		}
		c.Set(CtxAdminID, uint(adminID))
		c.Set(CtxRole, role)
		c.Next()
	}
}
