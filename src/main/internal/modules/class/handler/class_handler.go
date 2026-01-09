package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	service *service.ClassService
}

// NewClassHandler creates handler with injected service
func NewClassHandler(svc *service.ClassService) *ClassHandler {
	return &ClassHandler{service: svc}
}

// GetAllClass godoc
// @Summary Get all classes
// @Description Return Class list data
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Number of items per page (default: 10)"
// @Success 200 {object} response.BaseResponse
// @Router /academic/class/all [get]
func (h *ClassHandler) GetAllClass(c *gin.Context) {
	var req pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&req)

	data := h.service.FindAllClass(c, req)
	response.OK(c, "Success get all class", data)
}

// GetDetailClass godoc
// @Summary Get class detail
// @Description Return detail information class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class ID"
// @Success 200 {object} response.BaseResponse{data=dto.DetailClassResponse}
// @Router /academic/class/detail/{id} [get]
func (h *ClassHandler) GetDetailClass(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	detail := h.service.GetDetailClass(c, uint(id))
	response.OK(c, "Success get detail class", detail)
}

// CreateNewClass godoc
// @Summary Create a new class
// @Description Create New Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ModifyClassRequest true "Class data"
// @Success 201 {object} response.BaseResponse{data=dto.DetailClassResponse}
// @Router /academic/class/create [post]
func (h *ClassHandler) CreateNewClass(c *gin.Context) {
	var req dto.ModifyClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.CreateNewClass(c, req)
	observer.Trigger(model.EventClassChanged)
	response.OK(c, "Success create new class", resp)
}

// UpdateClass godoc
// @Summary Update a class
// @Description Update Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class ID"
// @Param request body dto.ModifyClassRequest true "Class data"
// @Success 200 {object} response.BaseResponse{data=dto.DetailClassResponse}
// @Router /academic/class/update/{id} [put]
func (h *ClassHandler) UpdateClass(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	var req dto.ModifyClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resp := h.service.ModifyClass(c, uint(id), req)
	observer.Trigger(model.EventClassChanged)
	response.OK(c, "Success update class", resp)
}

// DeleteClass godoc
// @Summary Delete a class
// @Description Delete Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Class ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/class/delete/{id} [delete]
func (h *ClassHandler) DeleteClass(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	h.service.DeleteClass(c, uint(id))
	observer.Trigger(model.EventClassChanged)
	response.OK(c, "Success delete class", gin.H{"id": id})
}

// MemberOfClass godoc
// @Summary Get class members
// @Description Get Member of Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param classId path int true "Class ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/class/{classId}/member [get]
func (h *ClassHandler) MemberOfClass(c *gin.Context) {
	idParam := c.Param("classId")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid Class ID", nil)
		return
	}

	res := h.service.MemberOfClass(uint(id))
	response.OK(c, "Success get class members", res)
}

// AddMemberOfClass godoc
// @Summary Add members to class
// @Description Add Member of Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ModifyClassMemberRequest true "Member data"
// @Success 200 {object} response.BaseResponse{data=dto.ModifyClassMemberRequest}
// @Router /academic/class/member/add [post]
func (h *ClassHandler) AddMemberOfClass(c *gin.Context) {
	var req dto.ModifyClassMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.AddMemberOfClass(req)
	observer.Trigger(model.EventClassChanged)
	response.OK(c, "Success add members", result)
}

// DeleteMemberOfClass godoc
// @Summary Remove member from class
// @Description Remove Member of Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Member ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/class/member/{id}/delete [delete]
func (h *ClassHandler) DeleteMemberOfClass(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequestError(c, "Invalid ID", nil)
		return
	}

	h.service.RemoveMemberOfClass(uint(id))
	observer.Trigger(model.EventClassChanged)
	response.OK(c, "Success remove member", gin.H{"id": id})
}

// DeleteMemberBatchOfClass godoc
// @Summary Batch remove members from class
// @Description Batch delete Member of Class
// @Tags Class
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ModifyClassMemberRequest true "Member data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/class/batch/delete [delete]
func (h *ClassHandler) DeleteMemberBatchOfClass(c *gin.Context) {
	var req dto.ModifyClassMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result := h.service.DeleteBatchMemberOfClass(req)
	observer.Trigger(model.EventClassChanged)
	response.OK(c, "Success batch remove members", result)
}
