package handler

import (
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/bucket"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/service"
	studentService "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/student/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	request2 "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/auth_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/gin-gonic/gin"
)

// CBTHandler handles CBT (student-facing) HTTP requests
type CBTHandler struct {
	examSessionService *service.ExamSessionService
	studentService     *studentService.StudentService
}

// NewCBTHandler creates handler with injected services
func NewCBTHandler(examSessionSvc *service.ExamSessionService, studentSvc *studentService.StudentService) *CBTHandler {
	return &CBTHandler{
		examSessionService: examSessionSvc,
		studentService:     studentSvc,
	}
}

// AuthCBTLogin godoc
// @Summary Student CBT login
// @Description Auth login for student
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth_request.CBTAuthRequest true "Login request"
// @Success 200 {object} response.BaseResponse
// @Router /auth/cbt/login [post]
func (h *CBTHandler) AuthCBTLogin(c *gin.Context) {
	var request auth_request.CBTAuthRequest
	_ = c.BindJSON(&request)

	resp := h.studentService.LoginByNISN(request)
	response.OK(c, "Success login", resp)
}

// ValidateToken godoc
// @Summary Validate CBT token
// @Description Validate exam session token
// @Tags CBT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionStartDoWork true "Token request"
// @Success 200 {object} response.BaseResponse
// @Router /cbt/token/validate [post]
func (h *CBTHandler) ValidateToken(c *gin.Context) {
	var request exam_request.ExamSessionStartDoWork
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	resp := h.examSessionService.ValidateTokenDo(claims, request)
	response.OK(c, "Success validate token", resp)
}

// RetrieveDetailSessionCbt godoc
// @Summary Retrieve CBT session
// @Description Retrieve detail of CBT session
// @Tags CBT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth_request.CBTSelectedSession true "Session request"
// @Success 200 {object} response.BaseResponse
// @Router /cbt/retrieve-session [post]
func (h *CBTHandler) RetrieveDetailSessionCbt(c *gin.Context) {
	var request auth_request.CBTSelectedSession
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	resp := h.studentService.RetrieveDetailSession(claims, request)
	response.OK(c, "Success retrieve token", resp)
}

// SubmitExamSession godoc
// @Summary Submit exam
// @Description Submit exam session
// @Tags CBT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionSubmit true "Submit request"
// @Success 200 {object} response.BaseResponse
// @Router /cbt/exam/submit [post]
func (h *CBTHandler) SubmitExamSession(c *gin.Context) {
	var request exam_request.ExamSessionSubmit
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	resp := h.examSessionService.SubmitExamSession(claims, request)
	response.OK(c, "Success submit token", resp)
}

// SyncAnswer godoc
// @Summary Sync answer
// @Description Sync student answer
// @Tags CBT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.ExamSessionSubmit true "Sync request"
// @Success 200 {object} response.BaseResponse
// @Router /cbt/exam/answer/sync [post]
func (h *CBTHandler) SyncAnswer(c *gin.Context) {
	var request exam_request.ExamSessionSubmit
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	h.examSessionService.SyncAnswer(claims, request)
	response.OK(c, "Success sync", nil)
}

// LastAnswer godoc
// @Summary Get last answer
// @Description Retrieve latest answer
// @Tags CBT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId query string true "Session ID"
// @Success 200 {object} response.BaseResponse
// @Router /cbt/exam/answer/latest [post]
func (h *CBTHandler) LastAnswer(c *gin.Context) {
	claims := jwt.GetDataClaims(c)
	resp := h.examSessionService.RetrieveLatestAnswer(claims, c.Query("sessionId"))
	response.OK(c, "Success sync", resp)
}

// ExamSessionSuspiciousActivityReport godoc
// @Summary Report suspicious activity
// @Description Report suspicious activity during exam
// @Tags CBT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.SuspiciousActivityReport true "Report request"
// @Success 200 {object} response.BaseResponse
// @Router /cbt/exam/report [post]
func (h *CBTHandler) ExamSessionSuspiciousActivityReport(c *gin.Context) {
	var request exam_request.SuspiciousActivityReport
	_ = c.BindJSON(&request)

	claims := jwt.GetDataClaims(c)
	h.examSessionService.SuspiciousActivityReport(claims, request)
	response.OK(c, "Success submit report", request)
}

// SessionInfo godoc
// @Summary Get session info
// @Description Get session info
// @Tags CBT
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body exam_request.SuspiciousActivityReport true "Request"
// @Success 200 {object} response.BaseResponse
// @Router /cbt/session [post]
func (h *CBTHandler) SessionInfo(c *gin.Context) {
	var request exam_request.SuspiciousActivityReport
	_ = c.BindJSON(&request)

	res := h.examSessionService.SessionInfo(request)
	response.OK(c, "Success session", res)
}

// UploadBase64 handles base64 file upload
func (h *CBTHandler) UploadBase64(c *gin.Context) {
	var request request2.UploadBase64Request
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(err)
	}

	minioCof := bucket.NewMinio()
	info, url := minioCof.UploadViaBase64(request.FileData, time.Now().Format("20060102"))

	response.OK(c, "Success", map[string]interface{}{
		"info": info,
		"url":  url,
	})
}
