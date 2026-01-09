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

// SessionHandler handles exam session-related HTTP requests
type SessionHandler struct {
	examSessionService *service.ExamSessionService
}

// NewSessionHandler creates handler with injected service
func NewSessionHandler(svc *service.ExamSessionService) *SessionHandler {
	return &SessionHandler{examSessionService: svc}
}

// GetAllExamSession godoc
// @Summary Get all exam sessions
// @Description Return paginated exam session list
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Items per page (default: 10)"
// @Success 200 {object} response.BaseResponse{data=database.Paginator}
// @Router /academic/exam/session/all [get]
func (h *SessionHandler) GetAllExamSession(c *gin.Context) {
	var request pagination.Request[map[string]interface{}]
	_ = c.BindQuery(&request)

	resp := h.examSessionService.GetAllExamSession(c, request)
	response.OK(c, "Success get all exam session", resp)
}

// GetExamSession godoc
// @Summary Get exam session detail
// @Description Return exam session detail by ID
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exam Session ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/detail/{id} [get]
func (h *SessionHandler) GetExamSession(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	detail := h.examSessionService.GetDetailExamSession(uint(id))
	response.OK(c, "Success get exam session", detail)
}

// CreateExamSession godoc
// @Summary Create exam session
// @Description Create a new exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ModifyExamSessionRequest true "Session data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/create [post]
func (h *SessionHandler) CreateExamSession(c *gin.Context) {
	var request exam_request.ModifyExamSessionRequest
	_ = c.BindJSON(&request)

	resp := h.examSessionService.CreateExamSession(c, request)
	observer.Trigger(model.EventExamSessionChanged)
	response.OK(c, "Success create exam session", resp)
}

// UpdateExamSession godoc
// @Summary Update exam session
// @Description Update an exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exam Session ID"
// @Param request body exam_request.ModifyExamSessionRequest true "Session data"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/update/{id} [put]
func (h *SessionHandler) UpdateExamSession(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var request exam_request.ModifyExamSessionRequest
	_ = c.BindJSON(&request)

	resp := h.examSessionService.UpdateExamSession(c, uint(id), request)
	observer.Trigger(model.EventExamSessionChanged)
	response.OK(c, "Success update exam session", resp)
}

// DeleteExamSession godoc
// @Summary Delete exam session
// @Description Delete an exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exam Session ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/delete/{id} [delete]
func (h *SessionHandler) DeleteExamSession(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.examSessionService.DeleteExamSession(uint(id))
	observer.Trigger(model.EventExamSessionChanged)
	response.OK(c, "Success delete exam session", gin.H{})
}

// GetAttendance godoc
// @Summary Get attendance
// @Description Get attendance for exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param exam_session_id query string true "Exam Session ID"
// @Param class_id query string false "Class ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/attendance [get]
func (h *SessionHandler) GetAttendance(c *gin.Context) {
	var request exam_request.ExamSessionAttendanceRequest
	_ = c.BindQuery(&request)

	resp := h.examSessionService.GetAllAttendance(request)
	response.OK(c, "Success get attendance", resp)
}

// DownloadAttendance handles attendance download
func (h *SessionHandler) DownloadAttendance(c *gin.Context) {
	var request exam_request.ExamSessionAttendanceRequest
	_ = c.BindQuery(&request)
	data := h.examSessionService.GetAllAttendance(request)
	h.examSessionService.ExportExamSessionAttendanceToExcel(c, data, request)
}

// ExamSessionMember godoc
// @Summary Get session members
// @Description Get members of an exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path string true "Session ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/member/{sessionId} [get]
func (h *SessionHandler) ExamSessionMember(c *gin.Context) {
	sessionId := c.Param("sessionId")
	res := h.examSessionService.ExamSessionMember(sessionId)
	response.OK(c, "Success get exam session member", res)
}

// ExamSessionReport godoc
// @Summary Get session report
// @Description Get report for exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/report [get]
func (h *SessionHandler) ExamSessionReport(c *gin.Context) {
	var request exam_request.ExamSessionReportRequest
	_ = c.BindQuery(&request)

	res := h.examSessionService.GetAllReport(request)
	response.OK(c, "Success get exam session report", res)
}

// ExamSessionGenerateReport godoc
// @Summary Generate session report
// @Description Generate report for exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionGenerateReportRequest true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/generate/report [post]
func (h *SessionHandler) ExamSessionGenerateReport(c *gin.Context) {
	var request exam_request.ExamSessionGenerateReportRequest
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	h.examSessionService.GenerateReportSession(request.SessionId, claims.SchoolCode)
	response.OK(c, "Your request still process, Please check your request to page 'Laporan Nilai'", request)
}

// ExamSessionAnswerResultStudent godoc
// @Summary Get student answer result
// @Description Get student answer result for exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/answer/student [get]
func (h *SessionHandler) ExamSessionAnswerResultStudent(c *gin.Context) {
	var request exam_request.ExamSessionStudentAnswer
	_ = c.BindQuery(&request)

	res := h.examSessionService.GetAnswerStudent(request)
	response.OK(c, "Success get exam session answer student", res)
}

// ExamSessionAnswerStudentCorrection godoc
// @Summary Correct student answer
// @Description Correct student answer for exam session
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionStudentAnswer true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/correction/answer/student [post]
func (h *SessionHandler) ExamSessionAnswerStudentCorrection(c *gin.Context) {
	var request exam_request.ExamSessionStudentAnswer
	_ = c.BindJSON(&request)

	res := h.examSessionService.CorrectionAnswerStudent(request)
	response.OK(c, "Success correction exam session answer student", res)
}

// ExamSessionReset godoc
// @Summary Reset session
// @Description Reset exam session for student
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionResetRequest true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/reset [post]
func (h *SessionHandler) ExamSessionReset(c *gin.Context) {
	var request exam_request.ExamSessionResetRequest
	_ = c.BindJSON(&request)

	h.examSessionService.ResetSessionStudent(request)
	response.OK(c, "Success reset session student", request)
}

// ExamSessionCorrectionScore godoc
// @Summary Correct score
// @Description Correct student score
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionCorrectionRequest true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/change/score [post]
func (h *SessionHandler) ExamSessionCorrectionScore(c *gin.Context) {
	var request exam_request.ExamSessionCorrectionRequest
	_ = c.BindJSON(&request)

	h.examSessionService.CorrectionScoreStudent(c, request)
	response.OK(c, "Success change score student", request)
}

// ResetSuspicious godoc
// @Summary Reset suspicious activity
// @Description Reset suspicious activity flag
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.SuspiciousActivityReport true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/cheat/reset [post]
func (h *SessionHandler) ResetSuspicious(c *gin.Context) {
	var request exam_request.SuspiciousActivityReport
	_ = c.BindJSON(&request)

	h.examSessionService.ResetSuspiciousActivity(request)
	response.OK(c, "Success session", request)
}

// ExamSessionRecalculate handles recalculation
func (h *SessionHandler) ExamSessionRecalculate(c *gin.Context) {
	h.examSessionService.CorrectionScoreUserMoreThan100()
}

// GetAllExamSessionToken godoc
// @Summary Get all tokens
// @Description Get all exam session tokens
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param exam_session_id query string true "Exam Session ID"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/token/all [get]
func (h *SessionHandler) GetAllExamSessionToken(c *gin.Context) {
	var request exam_request.ExamSessionTokenFilterRequest
	_ = c.BindQuery(&request)

	resp := h.examSessionService.GetAllToken(c, request)
	response.OK(c, "Success get exam session token", resp)
}

// CreateExamSessionToken godoc
// @Summary Generate token
// @Description Generate exam session token
// @Tags Exam Session
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionGenerateToken true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /academic/exam/session/token/generate [post]
func (h *SessionHandler) CreateExamSessionToken(c *gin.Context) {
	var request exam_request.ExamSessionGenerateToken
	_ = c.BindJSON(&request)

	resp := h.examSessionService.GenerateToken(c, request)
	response.OK(c, "Success create exam session token", resp)
}
