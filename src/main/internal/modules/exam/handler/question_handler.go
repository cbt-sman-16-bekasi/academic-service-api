package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/gin-gonic/gin"
)

// QuestionHandler handles exam question-related HTTP requests
type QuestionHandler struct {
	examService *service.ExamService
}

// NewQuestionHandler creates handler with injected service
func NewQuestionHandler(svc *service.ExamService) *QuestionHandler {
	return &QuestionHandler{examService: svc}
}

// GetDetailExamQuestion godoc
// @Summary Get question detail
// @Description Return question detail by ID
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Question ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam-question/detail/{id} [get]
func (h *QuestionHandler) GetDetailExamQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	question := h.examService.GetDetailExamQuestion(uint(id))
	response.OK(c, "Success get exam question detail", question)
}

// CreateExamQuestion godoc
// @Summary Create question
// @Description Create a new exam question
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ModifyExamQuestionRequest true "Question data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam-question/create [post]
func (h *QuestionHandler) CreateExamQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	resp := h.examService.CreateExamQuestion(claims, request)
	response.OK(c, "Success create exam question", resp)
}

// UpdateExamQuestion godoc
// @Summary Update question
// @Description Update an exam question
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Question ID"
// @Param request body exam_request.ModifyExamQuestionRequest true "Question data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam-question/update/{id} [put]
func (h *QuestionHandler) UpdateExamQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)
	id, _ := strconv.Atoi(c.Param("id"))

	resp := h.examService.UpdateExamQuestion(uint(id), request)
	response.OK(c, "Success update exam question", resp)
}

// DeleteExamQuestion godoc
// @Summary Delete question
// @Description Delete an exam question
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Question ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam-question/delete/{id} [delete]
func (h *QuestionHandler) DeleteExamQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.examService.DeleteExamQuestion(uint(id))
	response.OK(c, "Success delete exam question", gin.H{})
}
