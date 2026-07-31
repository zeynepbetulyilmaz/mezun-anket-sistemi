package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mezun-anket-backend/internal/apperror"
	"mezun-anket-backend/internal/middleware"
	"mezun-anket-backend/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type tokenLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

// POST /api/v1/auth/token-login
func (h *AuthHandler) TokenLogin(c *gin.Context) {
	var req tokenLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.Validation("Token alanı zorunludur."))
		return
	}

	jwtStr, grad, err := h.auth.LoginWithToken(req.Token)
	if err != nil {
		middleware.Fail(c, apperror.Unauthorized("Bağlantının süresi dolmuş veya geçersiz. Lütfen size gönderilen son e-postadaki linki kullanın."))
		return
	}

	c.JSON(http.StatusOK, middleware.SuccessEnvelope(gin.H{
		"accessToken": jwtStr,
		"graduate": gin.H{
			"firstName":      grad.FirstName,
			"departmentName": grad.DepartmentName,
			"graduationYear": grad.GraduationYear,
		},
	}))
}

type adminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// POST /api/v1/admin/login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.Validation("Kullanıcı adı ve şifre zorunludur."))
		return
	}
	jwtStr, admin, err := h.auth.AdminLogin(req.Username, req.Password)
	if err != nil {
		middleware.Fail(c, apperror.Unauthorized("Kullanıcı adı veya şifre hatalı."))
		return
	}
	c.JSON(http.StatusOK, middleware.SuccessEnvelope(gin.H{
		"accessToken": jwtStr,
		"role":        admin.Role,
	}))
}
