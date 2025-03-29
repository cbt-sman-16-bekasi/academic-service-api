package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/student_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/student_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

type StudentController struct {
	studentService *student_service.StudentService
}

func NewStudentController() *StudentController {
	return &StudentController{
		studentService: student_service.NewStudentService(),
	}
}

// GetAllStudent Get data Student All
// @Summary This endpoint about list of Student
// @Description Return Student list data
// @Tags Student
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
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=student.StudentClass}} "Pagination data Student"
// @Router /academic/student/all [get]
func (s *StudentController) GetAllStudent(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := s.studentService.AllStudent(req)
	response.SuccessResponse("Success get all student", data).Json(c)
}

// GetStudent Get detail of Student
// @Summary This endpoint about detail of Student
// @Description Return detail information Student
// @Tags Student
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int false "Student ID"
//
// @Success 200 {object} response.BaseResponse{data=class_response.DetailClassResponse} "Detail of Student"
// @Router /academic/student/detail/{id} [get]
func (s *StudentController) GetStudent(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := s.studentService.DetailStudent(uint(id))
	response.SuccessResponse("Success get student", detail).Json(c)
}

// CreateStudent Create New Student
// @Summary This endpoint about create new Student
// @Description Create New Student
// @Tags Student
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body student_request.StudentModifyRequest true "Body Request of Create New Student"
//
// @Success 200 {object} response.BaseResponse{data=student_response.DetailStudentResponse} "Create Student response"
// @Router /academic/student/create [post]
func (s *StudentController) CreateStudent(c *gin.Context) {
	var req student_request.StudentModifyRequest
	_ = c.BindJSON(&req)

	resp := s.studentService.CreateStudent(req)
	response.SuccessResponse("Success create student", resp).Json(c)
}

// UpdateStudent Update Student
// @Summary This endpoint about Update Student
// @Description Update Student
// @Tags Student
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Student ID"
// @Param request body student_request.StudentModifyRequest true "Body Request of Update Student"
//
// @Success 200 {object} response.BaseResponse{data=student_response.DetailStudentResponse} "Update Student response"
// @Router /academic/student/update/{id} [put]
func (s *StudentController) UpdateStudent(c *gin.Context) {
	var req student_request.StudentModifyRequest
	_ = c.BindJSON(&req)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := s.studentService.UpdateStudent(uint(id), req)
	response.SuccessResponse("Success update student", resp).Json(c)
}

// DeleteStudent Delete New Student
// @Summary This endpoint about Delete Student
// @Description Delete Student
// @Tags Student
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Student ID"
//
// @Success 200 {object} response.BaseResponse "Delete Student response"
// @Router /academic/student/delete/{id} [delete]
func (s *StudentController) DeleteStudent(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.studentService.DeleteById(uint(id))
	response.SuccessResponse("Success delete student", gin.H{"id": id}).Json(c)
}
