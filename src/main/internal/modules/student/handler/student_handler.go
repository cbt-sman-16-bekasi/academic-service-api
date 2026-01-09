package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	service *service.StudentService
}

// NewStudentHandler creates handler with injected service
func NewStudentHandler(svc *service.StudentService) *StudentHandler {
	return &StudentHandler{service: svc}
}

// GetAllStudent godoc
// @Summary Get all students
// @Tags Student
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.BaseResponse
// @Router /academic/student/all [get]
func (h *StudentHandler) GetAllStudent(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := h.service.AllStudent(c, req)
	response.OK(c, "Success get all student", data)
}

// GetStudent godoc
// @Summary Get student detail
// @Tags Student
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 200 {object} response.BaseResponse{data=dto.DetailStudentResponse}
// @Router /academic/student/detail/{id} [get]
func (h *StudentHandler) GetStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	detail := h.service.DetailStudent(uint(id))
	response.OK(c, "Success get student", detail)
}

// CreateStudent godoc
// @Summary Create a new student
// @Tags Student
// @Security BearerAuth
// @Param request body dto.StudentModifyRequest true "Student data"
// @Success 201 {object} response.BaseResponse{data=dto.DetailStudentResponse}
// @Router /academic/student/create [post]
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req dto.StudentModifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.CreateStudent(c, req)
	observer.Trigger(model.EventStudentChanged)
	response.OK(c, "Success create student", resp)
}

// UpdateStudent godoc
// @Summary Update a student
// @Tags Student
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Param request body dto.StudentModifyRequest true "Student data"
// @Success 200 {object} response.BaseResponse{data=dto.DetailStudentResponse}
// @Router /academic/student/update/{id} [put]
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	var req dto.StudentModifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.UpdateStudent(uint(id), req)
	observer.Trigger(model.EventStudentChanged)
	response.OK(c, "Success update student", resp)
}

// DeleteStudent godoc
// @Summary Delete a student
// @Tags Student
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/student/delete/{id} [delete]
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	h.service.DeleteStudent(uint(id))
	observer.Trigger(model.EventStudentChanged)
	response.OK(c, "Success delete student", gin.H{"id": id})
}

// DownloadTemplate godoc
// @Summary Download student upload template
// @Tags Student
// @Security BearerAuth
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200
// @Router /academic/student/template/download [get]
func (h *StudentHandler) DownloadTemplate(c *gin.Context) {
	h.service.DownloadTemplateUpload(c)
}

// UploadStudent godoc
// @Summary Upload students from excel
// @Tags Student
// @Security BearerAuth
// @Accept multipart/form-data
// @Param file formData file true "Excel file"
// @Success 200 {object} response.BaseResponse
// @Router /academic/student/template/upload [post]
func (h *StudentHandler) UploadStudent(c *gin.Context) {
	h.service.UploadTemplate(c)
}
