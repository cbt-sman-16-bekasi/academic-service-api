package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/gin-gonic/gin"
)

// TypeExamHandler handles type exam-related HTTP requests
type TypeExamHandler struct {
	typeExamService *service.TypeExamService
}

// NewTypeExamHandler creates handler with injected service
func NewTypeExamHandler(svc *service.TypeExamService) *TypeExamHandler {
	return &TypeExamHandler{typeExamService: svc}
}

// GetAllTypeExam godoc
// @Summary Get all type exams
// @Description Return paginated type exam list
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Items per page (default: 10)"
// @Success 200 {object} response.BaseResponse{data=database.Paginator}
// @Router /academic/exam/type-exam/all [get]
func (h *TypeExamHandler) GetAllTypeExam(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := h.typeExamService.GetAll(c, request)
	response.OK(c, "Success get all type Exam", resp)
}

// GetDetailTypeExam godoc
// @Summary Get type exam detail
// @Description Return type exam detail by ID
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Type Exam ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/type-exam/detail/{id} [get]
func (h *TypeExamHandler) GetDetailTypeExam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := h.typeExamService.GetDetail(uint(id))
	response.OK(c, "Success get type Exam", res)
}

// CreateTypeExam godoc
// @Summary Create type exam
// @Description Create a new type exam
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ModifyTypeExamRequest true "Type exam data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/type-exam/create [post]
func (h *TypeExamHandler) CreateTypeExam(c *gin.Context) {
	var request exam_request.ModifyTypeExamRequest
	_ = c.BindJSON(&request)

	res := h.typeExamService.CreateTypeExam(c, request)
	observer.Trigger(model.EventTypeExamChanged)
	response.OK(c, "Success create new type Exam", res)
}

// ModifyTypeExam godoc
// @Summary Update type exam
// @Description Update a type exam
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Type Exam ID"
// @Param request body exam_request.ModifyTypeExamRequest true "Type exam data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/type-exam/update/{id} [put]
func (h *TypeExamHandler) ModifyTypeExam(c *gin.Context) {
	var request exam_request.ModifyTypeExamRequest
	_ = c.BindJSON(&request)
	id, _ := strconv.Atoi(c.Param("id"))

	res := h.typeExamService.ModifyTypeExam(uint(id), request)
	observer.Trigger(model.EventTypeExamChanged)
	response.OK(c, "Success update type Exam", res)
}

// DeleteTypeExam godoc
// @Summary Delete type exam
// @Description Delete a type exam
// @Tags Type Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Type Exam ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/type-exam/delete/{id} [delete]
func (h *TypeExamHandler) DeleteTypeExam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.typeExamService.DeleteTypeExam(uint(id))
	observer.Trigger(model.EventTypeExamChanged)
	response.OK(c, "Success delete type Exam", nil)
}
