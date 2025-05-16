package controllers

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/observer"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"strconv"
)

// GetAllExamSession Get All exam session
// @Summary This endpoint about list of exam session
// @Description Return Exam session list data
// @Tags Exam Session
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
// @Success 200 {object} response.BaseResponse{data=database.Paginator{records=exam_response.ExamSessionListResponse}} "Pagination data"
// @Router /academic/exam/session/all [get]
func (e *ExamController) GetAllExamSession(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := e.examSessionService.GetAllExamSession(c, request)
	response.SuccessResponse("Success get all exam session", resp).Json(c)
}

// GetExamSession Get All exam session detail
// @Summary This endpoint about list of exam session detail
// @Description Return Exam session detail data
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Exam Session ID"
//
// @Success 200 {object} response.BaseResponse{data=exam_response.ExamDetailSessionResponse} "Detail data"
// @Router /academic/exam/session/detail/{id} [get]
func (e *ExamController) GetExamSession(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	detail := e.examSessionService.GetDetailExamSession(uint(id))
	response.SuccessResponse("Success get exam session", detail).Json(c)
}

// DeleteExamSession Delete exam session
// @Summary This endpoint about delete of exam session detail
// @Description Action Delete Exam session detail data
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Exam Session ID"
//
// @Success 200 {object} response.BaseResponse "Detail data"
// @Router /academic/exam/session/delete/{id} [delete]
func (e *ExamController) DeleteExamSession(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	e.examSessionService.DeleteExamSession(uint(id))
	observer.Trigger(model.EventExamSessionChanged)
	response.SuccessResponse("Success delete exam session", gin.H{}).Json(c)
}

// UpdateExamSession Update exam session
// @Summary This endpoint about Update of exam session detail
// @Description Action Update Exam session detail data
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Exam Session ID"
// @Param request body exam_request.ModifyExamSessionRequest true "Request Body"
//
// @Success 200 {object} response.BaseResponse{data=exam_request.ModifyExamSessionRequest} "Detail data"
// @Router /academic/exam/session/update/{id} [get]
func (e *ExamController) UpdateExamSession(c *gin.Context) {
	var idParam = c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var request exam_request.ModifyExamSessionRequest
	_ = c.BindJSON(&request)
	resp := e.examSessionService.UpdateExamSession(c, uint(id), request)
	observer.Trigger(model.EventExamSessionChanged)
	response.SuccessResponse("Success update exam session", resp).Json(c)
}

// CreateExamSession Create exam session
// @Summary This endpoint about Create of exam session detail
// @Description Action Create Exam session detail data
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body exam_request.ModifyExamSessionRequest true "Request Body"
//
// @Success 200 {object} response.BaseResponse{data=exam_request.ModifyExamSessionRequest} "Detail data"
// @Router /academic/exam/session/create [post]
func (e *ExamController) CreateExamSession(c *gin.Context) {
	var request exam_request.ModifyExamSessionRequest
	_ = c.BindJSON(&request)

	resp := e.examSessionService.CreateExamSession(c, request)
	observer.Trigger(model.EventExamSessionChanged)
	response.SuccessResponse("Success create exam session", resp).Json(c)
}

// GetAttendance Get attendance
// @Summary This endpoint about Get attendance of exam session
// @Description Return Get attendance of exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param exam_session_id query string true "Exam Session ID"
// @Param class_id query string false "Class ID"
//
// @Success 200 {object} response.BaseResponse{data=[]exam_response.ExamSessionAttendanceResponse} "Detail data"
// @Router /academic/exam/session/attendance [get]
func (e *ExamController) GetAttendance(c *gin.Context) {
	var request exam_request.ExamSessionAttendanceRequest
	err := c.BindQuery(&request)
	if err != nil {
		panic(err)
	}

	resp := e.examSessionService.GetAllAttendance(request)
	response.SuccessResponse("Success get attendance", resp).Json(c)
}

// DownloadAttendance Download attendance
// @Summary This endpoint about Download attendance of exam session
// @Description Return Download attendance of exam session
// @Tags Exam Session
// @Accept json
// @Produce octet-stream
// @Security BearerAuth
//
// @Param exam_session_id query string true "Exam Session ID"
// @Param class_id query string false "Class ID"
//
// @Success 200 {file} file "File downloaded"
// @Router /academic/exam/session/attendance/download [get]
func (e *ExamController) DownloadAttendance(c *gin.Context) {
	var request exam_request.ExamSessionAttendanceRequest
	_ = c.BindQuery(&request)
	data := e.examSessionService.GetAllAttendance(request)
	e.examSessionService.ExportExamSessionAttendanceToExcel(c, data, request)
}

// GetAllExamSessionToken Get 100 latest token
// @Summary This endpoint about Get 100 latest token
// @Description Return Get 100 latest token
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param exam_session_id query string true "Exam Session ID"
//
// @Success 200 {object} response.BaseResponse{data=[]exam_response.ExamSessionTokenResponse} "Detail data"
// @Router /academic/exam/session/token/all [get]
func (e *ExamController) GetAllExamSessionToken(c *gin.Context) {
	var request exam_request.ExamSessionTokenFilterRequest
	_ = c.BindQuery(&request)

	resp := e.examSessionService.GetAllToken(c, request)
	response.SuccessResponse("Success get exam session token", resp).Json(c)
}

// CreateExamSessionToken Create exam session token
// @Summary This endpoint about Create exam session token
// @Description Action Create exam session token
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param request body exam_request.ExamSessionGenerateToken true "Request Body"
//
// @Success 200 {object} response.BaseResponse{data=school.TokenExamSession} "Detail data"
// @Router /academic/exam/session/token/generate [post]
func (e *ExamController) CreateExamSessionToken(c *gin.Context) {
	var request exam_request.ExamSessionGenerateToken
	_ = c.BindJSON(&request)

	resp := e.examSessionService.GenerateToken(c, request)
	response.SuccessResponse("Success create exam session token", resp).Json(c)
}

func (e *ExamController) ValidateToken(c *gin.Context) {
	var request exam_request.ExamSessionStartDoWork
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	resp := e.examSessionService.ValidateTokenDo(claims, request)
	response.SuccessResponse("Success validate token", resp).Json(c)
}

func (e *ExamController) SubmitExamSession(c *gin.Context) {
	var request exam_request.ExamSessionSubmit
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	resp := e.examSessionService.SubmitExamSession(claims, request)
	response.SuccessResponse("Success submit token", resp).Json(c)
}

func (e *ExamController) ExamSessionMember(c *gin.Context) {
	var sessionId = c.Param("sessionId")

	res := e.examSessionService.ExamSessionMember(sessionId)
	response.SuccessResponse("Success get exam session member", res).Json(c)
}

func (e *ExamController) ExamSessionReport(c *gin.Context) {
	var request exam_request.ExamSessionReportRequest
	_ = c.BindQuery(&request)

	res := e.examSessionService.GetAllReport(request)
	response.SuccessResponse("Success get exam session report", res).Json(c)
}

func (e *ExamController) ExamSessionAnswerResultStudent(c *gin.Context) {
	var request exam_request.ExamSessionStudentAnswer
	_ = c.BindQuery(&request)

	res := e.examSessionService.GetAnswerStudent(request)
	response.SuccessResponse("Success get exam session answer student", res).Json(c)
}

func (e *ExamController) ExamSessionAnswerStudentCorrection(c *gin.Context) {
	var request exam_request.ExamSessionStudentAnswer
	_ = c.BindJSON(&request)

	res := e.examSessionService.CorrectionAnswerStudent(request)
	response.SuccessResponse("Success correction exam session answer student", res).Json(c)
}

func (e *ExamController) ExamSessionRecalculate(c *gin.Context) {
	e.examSessionService.CorrectionScoreUserMoreThan100()
}

func (e *ExamController) ExamSessionGenerateReport(c *gin.Context) {
	var request exam_request.ExamSessionGenerateReportRequest
	_ = c.BindJSON(&request)

	go e.examSessionService.GenerateReportSession(request.SessionId)
	response.SuccessResponse("Your request still process, Please check your request to page 'Laporan Nilai'", request).Json(c)
}
