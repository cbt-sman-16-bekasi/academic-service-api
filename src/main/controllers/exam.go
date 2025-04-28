package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/exam_service"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

type ExamController struct {
	typeExamService    *exam_service.TypeExamService
	examService        *exam_service.ExamService
	examSessionService *exam_service.ExamSessionService
}

func NewExamController() *ExamController {
	return &ExamController{
		typeExamService:    exam_service.NewTypeExamService(),
		examService:        exam_service.NewExamService(),
		examSessionService: exam_service.NewExamSessionService(),
	}
}

// GetAllExam Get All exam
// @Summary This endpoint about list of exam
// @Description Return Exam list data
// @Tags Exam
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
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=exam_response.ExamListResponse}} "Pagination data"
// @Router /academic/exam/all [get]
func (e *ExamController) GetAllExam(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := e.examService.GetAllExam(request)
	response.SuccessResponse("Success get data exam", resp).Json(c)
}

// GetDetailExam Update of exam
// @Summary This endpoint about detail of exam
// @Description Return Exam detail data
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Exam ID"
//
// @Success 200 {object} response.BaseResponse{data=exam_response.ExamDetailResponse} "Detail data"
// @Router /academic/exam/detail/{id} [get]
func (e *ExamController) GetDetailExam(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := e.examService.GetDetailExam(uint(id))
	response.SuccessResponse("Success get exam detail", detail).Json(c)
}

// CreateExam Create new of exam
// @Summary This endpoint about create of exam
// @Description Action of creation exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body exam_request.ModifyExamRequest true "Request Body"
//
// @Success 200 {object} response.BaseResponse{data=school.Exam} "Detail data"
// @Router /academic/exam/create [post]
func (e *ExamController) CreateExam(c *gin.Context) {
	var request exam_request.ModifyExamRequest
	_ = c.BindJSON(&request)

	resp := e.examService.CreateNewExam(request)
	response.SuccessResponse("Success create exam", resp).Json(c)
}

// UpdateExam Update of exam
// @Summary This endpoint about update of exam
// @Description Action of modify exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Exam ID"
// @Param request body exam_request.ModifyExamRequest true "Request Body"
//
// @Success 200 {object} response.BaseResponse{data=school.Exam} "Detail data"
// @Router /academic/exam/update/{id} [put]
func (e *ExamController) UpdateExam(c *gin.Context) {
	var request exam_request.ModifyExamRequest
	_ = c.BindJSON(&request)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := e.examService.UpdateExam(uint(id), request)
	response.SuccessResponse("Success update exam", resp).Json(c)
}

// DeleteExam Delete of exam
// @Summary This endpoint about Delete of exam
// @Description Action of Delete exam
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Exam ID"
//
// @Success 200 {object} response.BaseResponse "Detail data"
// @Router /academic/exam/delete/{id} [delete]
func (e *ExamController) DeleteExam(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)
	e.examService.DeleteExam(uint(id))
	response.SuccessResponse("Success delete exam", gin.H{}).Json(c)
}

// GetAllExamQuestion Get all question by EXAM ID
// @Summary This endpoint about get of question
// @Description Return list question
// @Tags Exam
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param examId path int true "Exam ID"
//
// @Success 200 {object} response.BaseResponse{data=[]school.ExamQuestion} "List data"
// @Router /academic/exam/{examId}/question [get]
func (e *ExamController) GetAllExamQuestion(c *gin.Context) {
	var idParam = c.Param("examId")
	id, _ := strconv.Atoi(idParam)
	examQuestion := e.examService.GetAllExamQuestion(uint(id))
	response.SuccessResponse("Success get exam questions", examQuestion).Json(c)
}

// GetDetailExamQuestion Update question
// @Summary This endpoint about Update of question
// @Description Return detail question
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Question ID"
//
// @Success 200 {object} response.BaseResponse{data=exam_response.DetailExamQuestionResponse} "Detail data"
// @Router /academic/exam/question/detail/{id} [get]
func (e *ExamController) GetDetailExamQuestion(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)
	examQuestion := e.examService.GetDetailExamQuestion(uint(id))
	response.SuccessResponse("Success get exam question detail", examQuestion).Json(c)
}

// CreateExamQuestion Update question
// @Summary This endpoint about Update of question
// @Description Return detail question
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body exam_request.ModifyExamQuestionRequest true "Request"
//
// @Success 200 {object} response.BaseResponse{data=exam_request.ModifyExamQuestionRequest} "Detail data"
// @Router /academic/exam/question/create [post]
func (e *ExamController) CreateExamQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)

	resp := e.examService.CreateExamQuestion(request)
	response.SuccessResponse("Success create exam question", resp).Json(c)
}

// UpdateExamQuestion Update question
// @Summary This endpoint about Update of question
// @Description Return detail question
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Question ID}"
// @Param request body exam_request.ModifyExamQuestionRequest true "Request"
//
// @Success 200 {object} response.BaseResponse{data=exam_request.ModifyExamQuestionRequest} "Detail data"
// @Router /academic/exam/question/update/{id} [put]
func (e *ExamController) UpdateExamQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)

	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	resp := e.examService.UpdateExamQuestion(uint(id), request)
	response.SuccessResponse("Success update exam question", resp).Json(c)
}

// DeleteExamQuestion Delete question
// @Summary This endpoint about Delete of question
// @Description Return detail question
// @Tags Exam Question
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Question ID}"
//
// @Success 200 {object} response.BaseResponse "Detail data"
// @Router /academic/exam/question/delete/{id} [put]
func (e *ExamController) DeleteExamQuestion(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)
	e.examService.DeleteExamQuestion(uint(id))

	response.SuccessResponse("Success delete exam question", gin.H{}).Json(c)
}

func (e *ExamController) DownloadTemplateQuestion(c *gin.Context) {
	var idParam = c.Param("examId")
	id, _ := strconv.Atoi(idParam)
	e.examService.DownloadTemplateUploadQuestion(uint(id), c)
}

func (e *ExamController) UploadQuestion(c *gin.Context) {
	e.examService.UploadQuestion(c)
}

func (e *ExamController) UploadBankQuestion(c *gin.Context) {
	e.examService.UploadBankQuestion(c)
}

func (e *ExamController) GetAllBankQuestion(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := e.examService.GetAllBankQuestion(request)
	response.SuccessResponse("Success get data bank question", resp).Json(c)
}

func (e *ExamController) GetDetailMasterBankQuestion(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)
	detail := e.examService.GetDetailBankQuestion(uint(id))
	response.SuccessResponse("Success get exam detail", detail).Json(c)
}

func (e *ExamController) CreateMasterBankQuestion(c *gin.Context) {
	var request exam_request.ModifyMasterBankQuestionRequest
	_ = c.BindJSON(&request)

	res := e.examService.CreateMasterBankQuestion(request)
	response.SuccessResponse("Success create master bank question", res).Json(c)
}

func (e *ExamController) UpdateMasterBankQuestion(c *gin.Context) {
	var request exam_request.ModifyMasterBankQuestionRequest
	_ = c.BindJSON(&request)

	id, _ := strconv.Atoi(c.Param("id"))
	resp := e.examService.UpdateMasterBankQuestion(uint(id), request)
	response.SuccessResponse("Success update master bank question", resp).Json(c)
}

func (e *ExamController) DeleteMasterBankQuestion(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)
	e.examService.DeleteMasterBankQuestion(uint(id))
	response.SuccessResponse("Success delete master bank question", gin.H{}).Json(c)
}

func (e *ExamController) GetQuestionByBankQuestionCode(c *gin.Context) {
	var idParam = c.Param("code")

	res := e.examService.GetQuestionByMasterCode(idParam)
	response.SuccessResponse("Success get exam question by bank question", res).Json(c)
}

func (e *ExamController) GetQuestionByBankQuestion(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	res := e.examService.GetBankQuestionById(uint(id))
	response.SuccessResponse("Success get exam bank question", res).Json(c)
}

func (e *ExamController) CreateBankQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)
	resp := e.examService.CreateBankQuestion(request)
	response.SuccessResponse("Success create exam question", resp).Json(c)
}

func (e *ExamController) UpdateBankQuestion(c *gin.Context) {
	var request exam_request.ModifyExamQuestionRequest
	_ = c.BindJSON(&request)
	id, _ := strconv.Atoi(c.Param("id"))

	res := e.examService.UpdateBankQuestion(uint(id), request)
	response.SuccessResponse("Success update exam question", res).Json(c)
}

func (e *ExamController) DeleteBankQuestion(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	e.examService.DeleteBankQuestion(uint(id))
	response.SuccessResponse("Success delete exam question", gin.H{}).Json(c)
}

func (e *ExamController) GetExamMember(c *gin.Context) {
	var idParam = c.Param("examCode")

	res := e.examService.GetExamMember(idParam)
	response.SuccessResponse("Success get exam member", res).Json(c)
}
