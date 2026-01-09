package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/observer"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/gin-gonic/gin"
)

// BankHandler handles question bank-related HTTP requests
type BankHandler struct {
	examService *service.ExamService
}

// NewBankHandler creates handler with injected service
func NewBankHandler(svc *service.ExamService) *BankHandler {
	return &BankHandler{examService: svc}
}

// GetAllBankQuestion godoc
// @Summary Get all bank questions
// @Description Return paginated question bank list
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Items per page (default: 10)"
// @Success 200 {object} response.BaseResponse{data=database.Paginator}
// @Router /academic/bank/all [get]
func (h *BankHandler) GetAllBankQuestion(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := h.examService.GetAllBankQuestion(c, request)
	response.OK(c, "Success get data bank question", resp)
}

// GetDetailMasterBankQuestion godoc
// @Summary Get bank question detail
// @Description Return bank question detail by ID
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Bank Question ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/detail/{id} [get]
func (h *BankHandler) GetDetailMasterBankQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	detail := h.examService.GetDetailBankQuestion(uint(id))
	response.OK(c, "Success get exam detail", detail)
}

// GetDetailMasterBankQuestionSubject godoc
// @Summary Get bank question by subject
// @Description Return bank question by subject
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/detail/bank/subject [get]
func (h *BankHandler) GetDetailMasterBankQuestionSubject(c *gin.Context) {
	var request exam_request.ModifyMasterBankQuestionRequest
	_ = c.BindQuery(&request)

	res := h.examService.GetDetailBankQuestionBySubject(request)
	response.OK(c, "Success detail master bank question", res)
}

// CreateMasterBankQuestion godoc
// @Summary Create master bank question
// @Description Create a new master bank question
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ModifyMasterBankQuestionRequest true "Bank question data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/create [post]
func (h *BankHandler) CreateMasterBankQuestion(c *gin.Context) {
	var request exam_request.ModifyMasterBankQuestionRequest
	_ = c.BindJSON(&request)

	res := h.examService.CreateMasterBankQuestion(jwt.GetIDClaims(c), request)
	observer.Trigger(model.EventBankQuestionChanged)
	response.OK(c, "Success create master bank question", res)
}

// UpdateMasterBankQuestion godoc
// @Summary Update master bank question
// @Description Update a master bank question
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Bank Question ID"
// @Param request body exam_request.ModifyMasterBankQuestionRequest true "Bank question data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/update/{id} [put]
func (h *BankHandler) UpdateMasterBankQuestion(c *gin.Context) {
	var request exam_request.ModifyMasterBankQuestionRequest
	_ = c.BindJSON(&request)
	id, _ := strconv.Atoi(c.Param("id"))

	resp := h.examService.UpdateMasterBankQuestion(jwt.GetIDClaims(c), uint(id), request)
	observer.Trigger(model.EventBankQuestionChanged)
	response.OK(c, "Success update master bank question", resp)
}

// DeleteMasterBankQuestion godoc
// @Summary Delete master bank question
// @Description Delete a master bank question
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Bank Question ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/delete/{id} [delete]
func (h *BankHandler) DeleteMasterBankQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.examService.DeleteMasterBankQuestion(uint(id))
	observer.Trigger(model.EventBankQuestionChanged)
	response.OK(c, "Success delete master bank question", gin.H{})
}

// GetQuestionByBankQuestionCode godoc
// @Summary Get questions by bank code
// @Description Get questions from a specific bank
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Bank Code"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/question/{code} [get]
func (h *BankHandler) GetQuestionByBankQuestionCode(c *gin.Context) {
	code := c.Param("code")
	res := h.examService.GetQuestionByMasterCode(code)
	response.OK(c, "Success get exam question by bank question", res)
}

// GetQuestionByBankQuestion godoc
// @Summary Get bank question by ID
// @Description Get a bank question by ID
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Question ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/question/detail/{id} [get]
func (h *BankHandler) GetQuestionByBankQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	res := h.examService.GetBankQuestionById(uint(id))
	response.OK(c, "Success get exam bank question", res)
}

// CreateBankQuestion godoc
// @Summary Create bank question
// @Description Create a new bank question
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ModifyExamQuestionRequest true "Question data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/question/create [post]
func (h *BankHandler) CreateBankQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)
	resp := h.examService.CreateBankQuestion(request)
	response.OK(c, "Success create exam question", resp)
}

// UpdateBankQuestion godoc
// @Summary Update bank question
// @Description Update a bank question
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Question ID"
// @Param request body exam_request.ModifyExamQuestionRequest true "Question data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/question/update/{id} [put]
func (h *BankHandler) UpdateBankQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)
	id, _ := strconv.Atoi(c.Param("id"))

	res := h.examService.UpdateBankQuestion(uint(id), request)
	response.OK(c, "Success update exam question", res)
}

// DeleteBankQuestion godoc
// @Summary Delete bank question
// @Description Delete a bank question
// @Tags Bank Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Question ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/bank/question/delete/{id} [delete]
func (h *BankHandler) DeleteBankQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.examService.DeleteBankQuestion(uint(id))
	response.OK(c, "Success delete exam question", gin.H{})
}

// UploadBankQuestion handles bank question upload
func (h *BankHandler) UploadBankQuestion(c *gin.Context) {
	h.examService.UploadBankQuestion(c)
}
