package handler

import (
	"fmt"
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	service *service.SubjectService
}

// NewSubjectHandler creates handler with injected service
func NewSubjectHandler(svc *service.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: svc}
}

// GetAllSubject godoc
// @Summary Get all subjects
// @Tags Curriculum
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {object} response.BaseResponse
// @Router /curriculum/subject/all [get]
func (h *SubjectHandler) GetAllSubject(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := h.service.GetAllSubject(c, request)
	response.OK(c, "Success get data subject", resp)
}

// GetSubject godoc
// @Summary Get subject by ID
// @Tags Curriculum
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Success 200 {object} response.BaseResponse
// @Router /curriculum/subject/detail/{id} [get]
func (h *SubjectHandler) GetSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	resp := h.service.GetSubject(id)
	response.OK(c, "Success get data subject", resp)
}

// CreateSubject godoc
// @Summary Create a new subject
// @Tags Curriculum
// @Security BearerAuth
// @Param request body dto.SubjectRequest true "Subject data"
// @Success 201 {object} response.BaseResponse
// @Router /curriculum/subject/create [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var req dto.SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.CreateSubject(c, req)
	observer.Trigger(model.EventSubjectsChanged)
	response.OK(c, "Success create subject", resp)
}

// UpdateSubject godoc
// @Summary Update a subject
// @Tags Curriculum
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Param request body dto.SubjectRequest true "Subject data"
// @Success 200 {object} response.BaseResponse
// @Router /curriculum/subject/update/{id} [put]
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	var req dto.SubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.UpdateSubject(c, id, req)
	observer.Trigger(model.EventSubjectsChanged)
	response.OK(c, "Success update subject", resp)
}

// DeleteSubject godoc
// @Summary Delete a subject
// @Tags Curriculum
// @Security BearerAuth
// @Param id path int true "Subject ID"
// @Success 200 {object} response.BaseResponse
// @Router /curriculum/subject/delete/{id} [delete]
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	h.service.DeleteSubject(id)
	observer.Trigger(model.EventSubjectsChanged)
	response.OK(c, fmt.Sprintf("Success delete subject %d", id), nil)
}
