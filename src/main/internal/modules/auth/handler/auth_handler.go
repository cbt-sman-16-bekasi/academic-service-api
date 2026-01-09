package handler

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates handler with injected service
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

// AuthLogin godoc
// @Summary Admin/Teacher login
// @Description Auth login for admin/teacher
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AuthRequest true "Login credentials"
// @Success 200 {object} response.Response{data=dto.AuthResponse}
// @Router /academic/auth/login [post]
func (h *AuthHandler) AuthLogin(c *gin.Context) {
	var req dto.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.Login(req.Username, req.Password)
	response.OK(c, "Success login", resp)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change user password
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.Response
// @Router /academic/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	h.service.ChangePassword(c, req)
	response.OK(c, "Success change password", nil)
}

// ChangeProfile godoc
// @Summary Change profile
// @Description Change user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangeProfileRequest true "Profile change data"
// @Success 200 {object} response.Response
// @Router /academic/auth/change-profile [post]
func (h *AuthHandler) ChangeProfile(c *gin.Context) {
	var req dto.ChangeProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	h.service.ChangeProfile(c, req)
	response.OK(c, "Success change profile", nil)
}
