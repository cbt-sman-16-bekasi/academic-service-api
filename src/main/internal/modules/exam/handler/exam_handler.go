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

// ExamHandler handles exam-related HTTP requests
type ExamHandler struct {
	examService *service.ExamService
}

// NewExamHandler creates handler with injected service
func NewExamHandler(svc *service.ExamService) *ExamHandler {
	return &ExamHandler{examService: svc}
}

// GetAllExam godoc
// @Summary Get all exams
// @Description Return paginated exam list
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Items per page (default: 10)"
// @Success 200 {object} response.BaseResponse{data=database.Paginator}
// @Router /academic/exam/all [get]
func (h *ExamHandler) GetAllExam(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := h.examService.GetAllExam(c, request)
	response.OK(c, "Success get data exam", resp)
}

// GetDetailExam godoc
// @Summary Get exam detail
// @Description Return exam detail by ID
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exam ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/detail/{id} [get]
func (h *ExamHandler) GetDetailExam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	detail := h.examService.GetDetailExam(uint(id))
	response.OK(c, "Success get exam detail", detail)
}

// CreateExam godoc
// @Summary Create new exam
// @Description Create a new exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ModifyExamRequest true "Exam data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/create [post]
func (h *ExamHandler) CreateExam(c *gin.Context) {
	var request exam_request.ModifyExamRequest
	_ = c.BindJSON(&request)

	resp := h.examService.CreateNewExam(c, request)
	observer.Trigger(model.EventExamChanged)
	response.OK(c, "Success create exam", resp)
}

// UpdateExam godoc
// @Summary Update exam
// @Description Update an existing exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exam ID"
// @Param request body exam_request.ModifyExamRequest true "Exam data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/update/{id} [put]
func (h *ExamHandler) UpdateExam(c *gin.Context) {
	var request exam_request.ModifyExamRequest
	_ = c.BindJSON(&request)
	id, _ := strconv.Atoi(c.Param("id"))

	resp := h.examService.UpdateExam(c, uint(id), request)
	observer.Trigger(model.EventExamChanged)
	response.OK(c, "Success update exam", resp)
}

// DeleteExam godoc
// @Summary Delete exam
// @Description Delete an exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exam ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/delete/{id} [delete]
func (h *ExamHandler) DeleteExam(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.examService.DeleteExam(uint(id))
	observer.Trigger(model.EventExamChanged)
	response.OK(c, "Success delete exam", gin.H{})
}

// GetExamMember godoc
// @Summary Get exam members
// @Description Get members of an exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param examCode path string true "Exam Code"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/member/{examCode} [get]
func (h *ExamHandler) GetExamMember(c *gin.Context) {
	examCode := c.Param("examCode")
	res := h.examService.GetExamMember(examCode)
	response.OK(c, "Success get exam member", res)
}

// GetAllExamQuestion godoc
// @Summary Get exam questions
// @Description Get all questions for an exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param examId path int true "Exam ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/{examId}/question [get]
func (h *ExamHandler) GetAllExamQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("examId"))
	questions := h.examService.GetAllExamQuestion(uint(id))
	response.OK(c, "Success get exam questions", questions)
}

// AddExamQuestionFromBank godoc
// @Summary Add questions from bank
// @Description Add questions from question bank to exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.AddExamQuestionFromBank true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/question/bank/add [post]
func (h *ExamHandler) AddExamQuestionFromBank(c *gin.Context) {
	var request exam_request.AddExamQuestionFromBank
	_ = c.BindJSON(&request)

	resp := h.examService.AddQuestionFromBank(request)
	response.OK(c, "Success add exam question", resp)
}

// DownloadTemplateQuestion handles template download
func (h *ExamHandler) DownloadTemplateQuestion(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("examId"))
	h.examService.DownloadTemplateUploadQuestion(uint(id), c)
}

// UploadQuestion handles question upload
func (h *ExamHandler) UploadQuestion(c *gin.Context) {
	h.examService.UploadQuestion(c)
}
