package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/user_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/user_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
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

func (s *UserController) GetAllUser(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	res := s.userService.GetAllUser(req)
	response.SuccessResponse("Success get data user", res).Json(c)
}

func (s *UserController) GetDetailUser(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	res := s.userService.GetUserById(uint(id))
	response.SuccessResponse("Success get data user", res).Json(c)
}

func (s *UserController) UpdateUser(c *gin.Context) {
	var req user_request.UserUpdateRequest
	_ = c.BindJSON(&req)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.userService.EnhanceSimpleUser(uint(id), req)
	response.SuccessResponse("Success update user data", req).Json(c)
}

func (s *UserController) ResetPassword(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.userService.ResetPasswordUser(uint(id))
	response.SuccessResponse("Berhasil reset password. Password baru sama dengan Username", nil).Json(c)
}

func (s *UserController) CreateUser(c *gin.Context) {
	var req user_request.UserUpdateRequest
	_ = c.BindJSON(&req)

	s.userService.CreateUser(req)
	response.SuccessResponse("Success create user data", req).Json(c)
}
