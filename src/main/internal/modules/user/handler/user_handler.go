package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates handler with injected service
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

// GetAllRoles godoc
// @Summary Get all roles
// @Tags User
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse
// @Router /academic/user/roles [get]
func (h *UserHandler) GetAllRoles(c *gin.Context) {
	data := h.service.GetAllRole()
	response.OK(c, "Success get data roles", data)
}

// GetAllUser godoc
// @Summary Get all users
// @Tags User
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.BaseResponse
// @Router /academic/user/all [get]
func (h *UserHandler) GetAllUser(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	res := h.service.GetAllUser(c, req)
	response.OK(c, "Success get data user", res)
}

// GetDetailUser godoc
// @Summary Get user detail
// @Tags User
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/user/detail/{id} [get]
func (h *UserHandler) GetDetailUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	res := h.service.GetUserById(uint(id))
	response.OK(c, "Success get data user", res)
}

// UpdateUser godoc
// @Summary Update user
// @Tags User
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body dto.UserUpdateRequest true "User data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/user/update/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	h.service.EnhanceSimpleUser(uint(id), req)
	response.OK(c, "Success update user data", req)
}

// ResetPassword godoc
// @Summary Reset user password
// @Tags User
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/user/reset-password/{id} [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	h.service.ResetPasswordUser(uint(id))
	response.OK(c, "Berhasil reset password. Password baru sama dengan Username", nil)
}

// CreateUser godoc
// @Summary Create new user
// @Tags User
// @Security BearerAuth
// @Param request body dto.UserUpdateRequest true "User data"
// @Success 201 {object} response.BaseResponse
// @Router /academic/user/create [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	h.service.CreateUser(c, req)
	response.OK(c, "Success create user data", req)
}
