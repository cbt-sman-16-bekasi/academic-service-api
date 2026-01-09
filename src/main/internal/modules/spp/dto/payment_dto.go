package dto

import (
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/entity"
)

// CreatePaymentRequest for creating new payment
type CreatePaymentRequest struct {
	StudentID    uint      `json:"student_id" binding:"required"`
	PaymentType  string    `json:"payment_type" binding:"required"`
	Amount       float64   `json:"amount" binding:"required,gt=0"`
	DueDate      time.Time `json:"due_date" binding:"required"`
	AcademicYear string    `json:"academic_year" binding:"required"`
	Semester     int       `json:"semester" binding:"required,min=1,max=2"`
	Notes        string    `json:"notes"`
}

// RecordPaymentRequest for recording payment transaction
type RecordPaymentRequest struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	TransactionID string  `json:"transaction_id"`
	CreatedBy     uint    `json:"created_by"`
	Notes         string  `json:"notes"`
}

// PaymentResponse represents payment data in response
type PaymentResponse struct {
	ID           uint       `json:"id"`
	StudentID    uint       `json:"student_id"`
	PaymentType  string     `json:"payment_type"`
	Amount       float64    `json:"amount"`
	PaidAmount   float64    `json:"paid_amount"`
	Outstanding  float64    `json:"outstanding"`
	Status       string     `json:"status"`
	DueDate      time.Time  `json:"due_date"`
	PaidAt       *time.Time `json:"paid_at,omitempty"`
	AcademicYear string     `json:"academic_year"`
	Semester     int        `json:"semester"`
	Notes        string     `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// TransactionResponse represents transaction data in response
type TransactionResponse struct {
	ID            uint      `json:"id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id,omitempty"`
	PaidAt        time.Time `json:"paid_at"`
	Notes         string    `json:"notes,omitempty"`
}

// PaymentDetailResponse includes payment with transactions
type PaymentDetailResponse struct {
	Payment      PaymentResponse       `json:"payment"`
	Transactions []TransactionResponse `json:"transactions"`
}

// ToPaymentResponse converts entity to response
func ToPaymentResponse(p *entity.Payment) PaymentResponse {
	return PaymentResponse{
		ID:           p.ID,
		StudentID:    p.StudentID,
		PaymentType:  p.PaymentType,
		Amount:       p.Amount,
		PaidAmount:   p.PaidAmount,
		Outstanding:  p.Amount - p.PaidAmount,
		Status:       p.Status,
		DueDate:      p.DueDate,
		PaidAt:       p.PaidAt,
		AcademicYear: p.AcademicYear,
		Semester:     p.Semester,
		Notes:        p.Notes,
		CreatedAt:    p.CreatedAt,
	}
}

// ToTransactionResponse converts entity to response
func ToTransactionResponse(t *entity.PaymentTransaction) TransactionResponse {
	return TransactionResponse{
		ID:            t.ID,
		Amount:        t.Amount,
		PaymentMethod: t.PaymentMethod,
		TransactionID: t.TransactionID,
		PaidAt:        t.PaidAt,
		Notes:         t.Notes,
	}
}
