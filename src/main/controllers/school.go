package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/class_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/student_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/class_service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/school_service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/student_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

type SchoolController struct {
	srv            *school_service.SchoolService
	classService   *class_service.ClassService
	studentService *student_service.StudentService
}

func NewSchoolController() *SchoolController {
	return &SchoolController{
		srv:            school_service.NewSchoolService(),
		classService:   class_service.NewClassService(),
		studentService: student_service.NewStudentService(),
	}
}

func (s *SchoolController) GetSchool(c *gin.Context) {
	data := s.srv.RetrieveDetailSchool(c)
	response.SuccessResponse("Success get data school_repository", data).Json(c)
}

func (s *SchoolController) GetAllClassCode(c *gin.Context) {
	data := s.srv.GetAllClassCode()
	response.SuccessResponse("Success get all class code", data).Json(c)
}

func (s *SchoolController) GetAllSubject(c *gin.Context) {

}

func (s *SchoolController) GetAllClass(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := s.classService.FindAllClass(req)
	response.SuccessResponse("Success get all class", data).Json(c)
}

func (s *SchoolController) GetDetailClass(c *gin.Context) {

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := s.classService.GetDetailClass(uint(id))
	response.SuccessResponse("Success get detail class", detail).Json(c)
}

func (s *SchoolController) CreateNewClass(c *gin.Context) {
	var req class_request.ModifyClassRequest
	_ = c.BindJSON(&req)

	resp := s.classService.CreateNewClass(req)
	response.SuccessResponse("Success create new class", resp).Json(c)
}

func (s *SchoolController) UpdateClass(c *gin.Context) {
	var req class_request.ModifyClassRequest
	_ = c.BindJSON(&req)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := s.classService.ModifyClass(uint(id), req)
	response.SuccessResponse("Success update class", resp).Json(c)
}

func (s *SchoolController) DeleteClass(c *gin.Context) {

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.classService.DeleteById(uint(id))
	response.SuccessResponse("Success delete class", gin.H{"id": id}).Json(c)
}

func (s *SchoolController) GetAllClassSubject(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := s.srv.GetAllClassSubject(req)
	response.SuccessResponse("Success get all class subject", data).Json(c)
}

func (s *SchoolController) GetClassSubject(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := s.srv.GetDetailClassSubject(uint(id))
	response.SuccessResponse("Success get class subject", detail).Json(c)
}

func (s *SchoolController) CreateClassSubject(c *gin.Context) {
}

func (s *SchoolController) ModifyClassSubject(c *gin.Context) {
}

func (s *SchoolController) DeleteClassSubject(c *gin.Context) {

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.srv.DeleteClassSubject(uint(id))
	response.SuccessResponse("Success delete class subject", gin.H{"id": id}).Json(c)
}

func (s *SchoolController) GetAllStudent(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := s.studentService.AllStudent(req)
	response.SuccessResponse("Success get all student", data).Json(c)
}

func (s *SchoolController) GetStudent(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := s.studentService.DetailStudent(uint(id))
	response.SuccessResponse("Success get student", detail).Json(c)
}

func (s *SchoolController) CreateStudent(c *gin.Context) {
	var req student_request.StudentModifyRequest
	_ = c.BindJSON(&req)

	resp := s.studentService.CreateStudent(req)
	response.SuccessResponse("Success create student", resp).Json(c)
}

func (s *SchoolController) UpdateStudent(c *gin.Context) {
	var req student_request.StudentModifyRequest
	_ = c.BindJSON(&req)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := s.studentService.UpdateStudent(uint(id), req)
	response.SuccessResponse("Success update student", resp).Json(c)
}

func (s *SchoolController) DeleteStudent(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	s.studentService.DeleteById(uint(id))
	response.SuccessResponse("Success delete student", gin.H{"id": id}).Json(c)
}
