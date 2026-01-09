package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/gin-gonic/gin"
)

type TeacherHandler struct {
	service *service.TeacherService
}

// NewTeacherHandler creates handler with injected service
func NewTeacherHandler(svc *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{service: svc}
}

// GetAllTeacher godoc
// @Summary Get all teachers
// @Tags Teacher
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.BaseResponse
// @Router /academic/teacher/all [get]
func (h *TeacherHandler) GetAllTeacher(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	paging := h.service.GetAllTeacher(c, req)
	response.OK(c, "Success get all teacher", paging)
}

// GetTeacher godoc
// @Summary Get teacher detail
// @Tags Teacher
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Success 200 {object} response.BaseResponse{data=dto.TeacherDetailResponse}
// @Router /academic/teacher/detail/{id} [get]
func (h *TeacherHandler) GetTeacher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	resp := h.service.DetailTeacher(uint(id))
	response.OK(c, "Success get teacher", resp)
}

// CreateTeacher godoc
// @Summary Create a new teacher
// @Tags Teacher
// @Security BearerAuth
// @Param request body dto.TeacherModifyRequest true "Teacher data"
// @Success 201 {object} response.BaseResponse{data=dto.TeacherDetailResponse}
// @Router /academic/teacher/create [post]
func (h *TeacherHandler) CreateTeacher(c *gin.Context) {
	var req dto.TeacherModifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.CreateTeacher(c, req)
	observer.Trigger(model.EventTeacherChanged)
	response.OK(c, "Success create teacher", resp)
}

// UpdateTeacher godoc
// @Summary Update a teacher
// @Tags Teacher
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Param request body dto.TeacherModifyRequest true "Teacher data"
// @Success 200 {object} response.BaseResponse{data=dto.TeacherDetailResponse}
// @Router /academic/teacher/update/{id} [put]
func (h *TeacherHandler) UpdateTeacher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	var req dto.TeacherModifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.UpdateTeacher(c, uint(id), req)
	observer.Trigger(model.EventTeacherChanged)
	response.OK(c, "Success update teacher", resp)
}

// DeleteTeacher godoc
// @Summary Delete a teacher
// @Tags Teacher
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/teacher/delete/{id} [delete]
func (h *TeacherHandler) DeleteTeacher(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	h.service.DeleteTeacher(uint(id))
	observer.Trigger(model.EventTeacherChanged)
	response.OK(c, "Success delete teacher", gin.H{"id": id})
}

// GetTeacherSubjectClassList godoc
// @Summary Get teacher's class subjects
// @Tags Teacher
// @Security BearerAuth
// @Param teacherId path int true "Teacher ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/teacher/{teacherId}/subject-class [get]
func (h *TeacherHandler) GetTeacherSubjectClassList(c *gin.Context) {
	idParam := c.Param("teacherId")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid Teacher ID", nil)
		return
	}

	resp := h.service.GetAllSubjectClassTeacher(uint(id))
	response.OK(c, "Success get subject class list", resp)
}

// GetDetailTeacherSubject godoc
// @Summary Get detail of teacher's class subject
// @Tags Teacher
// @Security BearerAuth
// @Param id path int true "Class Subject ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/teacher/subject-class/detail/{id} [get]
func (h *TeacherHandler) GetDetailTeacherSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	resp := h.service.GetDetailSubjectClassTeacher(uint(id))
	response.OK(c, "Success get subject class", resp)
}

// CreateTeacherSubject godoc
// @Summary Create teacher's class subject assignment
// @Tags Teacher
// @Security BearerAuth
// @Param request body dto.TeacherMappingSubjectClass true "Class subject data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/teacher/subject-class/create [post]
func (h *TeacherHandler) CreateTeacherSubject(c *gin.Context) {
	var req dto.TeacherMappingSubjectClass
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.CreateSubjectClassTeacher(req)
	response.OK(c, "Success create subject class", resp)
}

// UpdateTeacherSubject godoc
// @Summary Update teacher's class subject assignment
// @Tags Teacher
// @Security BearerAuth
// @Param id path int true "Class Subject ID"
// @Param request body dto.TeacherMappingSubjectClass true "Class subject data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/teacher/subject-class/update/{id} [put]
func (h *TeacherHandler) UpdateTeacherSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	var req dto.TeacherMappingSubjectClass
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.UpdateSubjectClassTeacher(uint(id), req)
	response.OK(c, "Success update subject class", resp)
}

// DeleteTeacherSubject godoc
// @Summary Delete teacher's class subject assignment
// @Tags Teacher
// @Security BearerAuth
// @Param id path int true "Class Subject ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/teacher/subject-class/delete/{id} [delete]
func (h *TeacherHandler) DeleteTeacherSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	h.service.DeleteSubjectClassTeacher(uint(id))
	response.OK(c, "Success delete subject class", gin.H{"id": id})
}
