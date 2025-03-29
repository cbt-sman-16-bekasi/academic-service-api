package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/user_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/server/response"
)

type UserController struct {
	userService *user_service.UserService
}

func NewUserController() *UserController {
	return &UserController{userService: user_service.NewUserService()}
}

// GetAllRoles Get all Roles
// @Summary This endpoint about data list roles
// @Description Return list role
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Success 200 {object} response.BaseResponse{data=[]user.Role} "List of Role"
// @Router /academic/user/roles [get]
func (s *UserController) GetAllRoles(c *gin.Context) {
	data := s.userService.GetAllRole()
	response.SuccessResponse("Success get data roles", data).Json(c)
}
