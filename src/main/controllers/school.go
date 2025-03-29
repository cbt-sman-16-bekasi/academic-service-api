package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/auth_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/class_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/auth_service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/school_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

type SchoolController struct {
	srv         *school_service.SchoolService
	authService *auth_service.AuthService
}

func NewSchoolController() *SchoolController {
	return &SchoolController{
		srv:         school_service.NewSchoolService(),
		authService: auth_service.NewAuthService(),
	}
}

// GetSchool Get data detail about school
// @Summary This endpoint about school information
// @Description Return school data information
// @Tags School
// @Accept json
// @Produce json
// @Security
//
// @Param schoolCode query string true "Required school code"
// @Success 200 {object} response.BaseResponse{data=school_response.DetailSchool} "School information"
// @Router /academic/school [get]
func (s *SchoolController) GetSchool(c *gin.Context) {
	data := s.srv.RetrieveDetailSchool(c)
	response.SuccessResponse("Success get data school_repository", data).Json(c)
}

// GetAllClassCode Get data all class code
// @Summary This endpoint about Class Code like 10, 11, 12
// @Description Return Class Code list data
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Success 200 {object} response.BaseResponse{data=[]class_response.ClassCodeResponse} "Class Code Response"
// @Router /academic/class-code [get]
func (s *SchoolController) GetAllClassCode(c *gin.Context) {
	data := s.srv.GetAllClassCode()
	response.SuccessResponse("Success get all class code", data).Json(c)
}

// GetAllSubject Get all data subjects
// @Summary This endpoint get data all subjects like Matematika and others
// @Description Return Subject list data
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Success 200 {object} response.BaseResponse{data=[]curriculum.Subject} "subjects Response"
// @Router /academic/subjects [get]
func (s *SchoolController) GetAllSubject(c *gin.Context) {
	resp := s.srv.GetAllSubject()
	response.SuccessResponse("Success get all subject", resp).Json(c)
}

// GetAllClassSubject Get data Class Subject
// @Summary This endpoint about list of class Subject
// @Description Return Class Subject list data
// @Tags Class
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
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=school.ClassSubject}} "Pagination data class Subject"
// @Router /academic/class/subject/all [get]
func (s *SchoolController) GetAllClassSubject(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := s.srv.GetAllClassSubject(req)
	response.SuccessResponse("Success get all class subject", data).Json(c)
}

// GetClassSubject Get detail of Class
// @Summary This endpoint about detail of class
// @Description Return detail information class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int false "Class ID"
//
// @Success 200 {object} response.BaseResponse{data=class_response.DetailClassResponse} "Detail of class"
// @Router /academic/class/subject/detail/{id} [get]
func (s *SchoolController) GetClassSubject(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := s.srv.GetDetailClassSubject(uint(id))
	response.SuccessResponse("Success get class subject", detail).Json(c)
}

// CreateClassSubject Create new data class subject
// @Summary This endpoint Create new data class subject
// @Description Create new data class subject
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body class_request.ModifyClassSubject true "Body Request of Create new Class subject"
//
// @Success 200 {object} response.BaseResponse{data=class_response.DetailClassSubjectResponse} "Class subject Response"
// @Router /academic/class/subject/create [post]
func (s *SchoolController) CreateClassSubject(c *gin.Context) {
	var request class_request.ModifyClassSubject
	_ = c.BindJSON(&request)

	res := s.srv.CreateClassSubject(request)
	response.SuccessResponse("Success create class subject", res).Json(c)
}

// ModifyClassSubject Update new data class subject
// @Summary This endpoint Update new data class subject
// @Description Update new data class subject
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body class_request.ModifyClassSubject true "Body Request of Update new Class subject"
// @Param request path int true "Class Subject ID"
//
// @Success 200 {object} response.BaseResponse{data=class_response.DetailClassSubjectResponse} "Class subject Response"
// @Router /academic/class/subject/update/{id} [put]
func (s *SchoolController) ModifyClassSubject(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var request class_request.ModifyClassSubject
	_ = c.BindJSON(&request)

	res := s.srv.UpdateClassSubject(uint(id), request)
	response.SuccessResponse("Success update class subject", res).Json(c)
}

// DeleteClassSubject Delete New ClassSubject
// @Summary This endpoint about Delete ClassSubject
// @Description Delete ClassSubject
// @Tags Student
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "ClassSubject ID"
//
// @Success 200 {object} response.BaseResponse "Delete ClassSubject response"
// @Router /academic/class/subject/delete/{id} [delete]
func (s *SchoolController) DeleteClassSubject(c *gin.Context) {

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.srv.DeleteClassSubject(uint(id))
	response.SuccessResponse("Success delete class subject", gin.H{"id": id}).Json(c)
}

// AuthLogin Admin auth login
// @Summary This endpoint about auth login
// @Description Auth login for admin
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body auth_request.AuthRequest true "Body Request of Login"
//
// @Success 200 {object} response.BaseResponse{data=auth_response.AuthResponse} "Login Response"
// @Router /academic/auth/login [post]
func (s *SchoolController) AuthLogin(c *gin.Context) {
	var request auth_request.AuthRequest
	_ = c.BindJSON(&request)

	resp := s.authService.Login(request.Username, request.Password)
	response.SuccessResponse("Success login", resp).Json(c)
}

// GetDashboard Get data dashboard user
// @Summary This endpoint about dashboard user
// @Description Return dashboard user
// @Tags School
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Success 200 {object} response.BaseResponse{data=school_response.DashboardResponse} "Dashboard response"
// @Router /academic/dashboard [get]
func (s *SchoolController) GetDashboard(c *gin.Context) {
	dt := s.srv.DashboardUser()
	response.SuccessResponse("Success get dashboard", dt).Json(c)
}
