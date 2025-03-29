package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/teacher_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/teacher_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

type TeacherController struct {
	teacherService *teacher_service.TeacherService
}

func NewTeacherController() *TeacherController {
	return &TeacherController{
		teacherService: teacher_service.NewTeacherService(),
	}
}

// GetAllTeacher Get data Teacher All
// @Summary This endpoint about list of Teacher
// @Description Return Teacher list data
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param page query int false "Page number (default: 1)" default(1)
// @Param size query int false "Number of items per page (default: 10)" default(10)
// @Param search.key query string false "Search key field (e.g., name, email, etc.)"
// @Param search.value query string false "Search value for the specified key"
// @Param filter query object false "Filter conditions (varies depending on the API)"
//
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=teacher.Teacher}} "Pagination data Teacher"
// @Router /academic/teacher/all [get]
func (s *TeacherController) GetAllTeacher(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)
	paging := s.teacherService.GetAllTeacher(req)
	response.SuccessResponse("Success get all teacher", paging).Json(c)
}

// GetTeacher Get detail of Teacher
// @Summary This endpoint about detail of Teacher
// @Description Return detail information Teacher
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int false "Teacher ID"
//
// @Success 200 {object} response.BaseResponse{data=teacher_response.TeacherDetailResponse} "Detail of Teacher"
// @Router /academic/teacher/detail/{id} [get]
func (s *TeacherController) GetTeacher(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := s.teacherService.DetailTeacher(uint(id))
	response.SuccessResponse("Success get teacher", resp).Json(c)
}

// CreateTeacher Create New Teacher
// @Summary This endpoint about create new Teacher
// @Description Create New Teacher
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body teacher_request.TeacherModifyRequest true "Body Request of Create New Teacher"
//
// @Success 200 {object} response.BaseResponse{data=teacher_response.TeacherDetailResponse} "Create Teacher response"
// @Router /academic/teacher/create [post]
func (s *TeacherController) CreateTeacher(c *gin.Context) {
	var req teacher_request.TeacherModifyRequest
	_ = c.BindJSON(&req)
	resp := s.teacherService.CreateTeacher(req)
	response.SuccessResponse("Success create teacher", resp).Json(c)
}

// UpdateTeacher Update Teacher
// @Summary This endpoint about Update Teacher
// @Description Update Teacher
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Teacher ID"
// @Param request body teacher_request.TeacherModifyRequest true "Body Request of Update Teacher"
//
// @Success 200 {object} response.BaseResponse{data=teacher_response.TeacherDetailResponse} "Update Teacher response"
// @Router /academic/teacher/update/{id} [put]
func (s *TeacherController) UpdateTeacher(c *gin.Context) {
	var req teacher_request.TeacherModifyRequest
	_ = c.BindJSON(&req)
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)
	resp := s.teacherService.UpdateTeacher(uint(id), req)
	response.SuccessResponse("Success update teacher", resp).Json(c)
}

// DeleteTeacher Delete New Teacher
// @Summary This endpoint about Delete Teacher
// @Description Delete Teacher
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Teacher ID"
//
// @Success 200 {object} response.BaseResponse "Delete Teacher response"
// @Router /academic/teacher/delete/{id} [delete]
func (s *TeacherController) DeleteTeacher(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)
	s.teacherService.DeleteById(uint(id))
	response.SuccessResponse("Success delete teacher", gin.H{"id": id}).Json(c)
}
