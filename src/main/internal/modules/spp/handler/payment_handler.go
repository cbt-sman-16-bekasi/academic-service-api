package handler

import (
	"strconv"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service *service.PaymentService
}

// NewPaymentHandler creates handler with injected service
func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: svc}
}

// GetStudentPayments godoc
// @Summary Get student payments
// @Tags SPP
// @Security BearerAuth
// @Param student_id query int true "Student ID"
// @Success 200 {object} response.BaseResponse
// @Router /spp/student/payments [get]
func (h *PaymentHandler) GetStudentPayments(c *gin.Context) {
	schoolCode := jwt.GetSchoolCode(c)
	studentID, err := strconv.ParseUint(c.Query("student_id"), 10, 32)
	if err != nil {
		response.BadRequestError(c, "Invalid student_id", nil)
		return
	}

	payments, err := h.service.GetStudentPayments(schoolCode, uint(studentID))
	if err != nil {
		response.InternalError(c, err.Error(), nil)
		return
	}

	response.OK(c, "Success get student payments", payments)
}

// GetPendingPayments godoc
// @Summary Get all pending payments
// @Tags SPP
// @Security BearerAuth
// @Success 200 {object} response.BaseResponse
// @Router /spp/pending [get]
func (h *PaymentHandler) GetPendingPayments(c *gin.Context) {
	schoolCode := jwt.GetSchoolCode(c)

	payments, err := h.service.GetPendingPayments(schoolCode)
	if err != nil {
		response.InternalError(c, err.Error(), nil)
		return
	}

	response.OK(c, "Success get pending payments", payments)
}

// CreatePayment godoc
// @Summary Create new payment
// @Tags SPP
// @Security BearerAuth
// @Param request body dto.CreatePaymentRequest true "Payment data"
// @Success 201 {object} response.BaseResponse
// @Router /spp/create [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	schoolCode := jwt.GetSchoolCode(c)

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.service.CreatePayment(schoolCode, req)
	if err != nil {
		response.InternalError(c, err.Error(), nil)
		return
	}

	response.OK(c, "Payment created successfully", result)
}

// RecordPayment godoc
// @Summary Record payment transaction
// @Tags SPP
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Param request body dto.RecordPaymentRequest true "Transaction data"
// @Success 200 {object} response.BaseResponse
// @Router /spp/{id}/pay [post]
func (h *PaymentHandler) RecordPayment(c *gin.Context) {
	schoolCode := jwt.GetSchoolCode(c)
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequestError(c, "Invalid payment id", nil)
		return
	}

	var req dto.RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// Set created_by from JWT
	claims := jwt.GetDataClaims(c)
	req.CreatedBy = uint(claims.Id)

	if err := h.service.RecordPayment(schoolCode, uint(paymentID), req); err != nil {
		response.InternalError(c, err.Error(), nil)
		return
	}

	response.OK(c, "Payment recorded successfully", nil)
}

// GetPaymentDetail godoc
// @Summary Get payment detail with transactions
// @Tags SPP
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Success 200 {object} response.BaseResponse
// @Router /spp/{id}/detail [get]
func (h *PaymentHandler) GetPaymentDetail(c *gin.Context) {
	schoolCode := jwt.GetSchoolCode(c)
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequestError(c, "Invalid payment id", nil)
		return
	}

	result, err := h.service.GetPaymentDetail(schoolCode, uint(paymentID))
	if err != nil {
		response.NotFoundError(c, err.Error())
		return
	}

	response.OK(c, "Success get payment detail", result)
}
