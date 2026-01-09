package service

import (
	"errors"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/entity"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/repository"
)

type PaymentService struct {
	repo *repository.PaymentRepository
}

// NewPaymentService creates service with injected repository
func NewPaymentService(repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

// GetStudentPayments gets all payments for a student
func (s *PaymentService) GetStudentPayments(schoolCode string, studentID uint) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.FindByStudent(schoolCode, studentID)
	if err != nil {
		return nil, err
	}

	var responses []dto.PaymentResponse
	for _, p := range payments {
		responses = append(responses, dto.ToPaymentResponse(&p))
	}
	return responses, nil
}

// GetPendingPayments gets all pending payments for a school
func (s *PaymentService) GetPendingPayments(schoolCode string) ([]dto.PaymentResponse, error) {
	payments, err := s.repo.FindPending(schoolCode)
	if err != nil {
		return nil, err
	}

	var responses []dto.PaymentResponse
	for _, p := range payments {
		responses = append(responses, dto.ToPaymentResponse(&p))
	}
	return responses, nil
}

// CreatePayment creates a new payment record
func (s *PaymentService) CreatePayment(schoolCode string, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	payment := &entity.Payment{
		SchoolCode:   schoolCode,
		StudentID:    req.StudentID,
		PaymentType:  req.PaymentType,
		Amount:       req.Amount,
		DueDate:      req.DueDate,
		AcademicYear: req.AcademicYear,
		Semester:     req.Semester,
		Notes:        req.Notes,
		Status:       "pending",
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	response := dto.ToPaymentResponse(payment)
	return &response, nil
}

// RecordPayment records a payment transaction
func (s *PaymentService) RecordPayment(schoolCode string, paymentID uint, req dto.RecordPaymentRequest) error {
	payment, err := s.repo.FindByID(schoolCode, paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Create transaction
	now := time.Now()
	tx := &entity.PaymentTransaction{
		SchoolCode:    schoolCode,
		PaymentID:     paymentID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		TransactionID: req.TransactionID,
		PaidAt:        now,
		CreatedBy:     req.CreatedBy,
		Notes:         req.Notes,
	}

	if err := s.repo.CreateTransaction(tx); err != nil {
		return err
	}

	// Update payment status
	payment.PaidAmount += req.Amount
	if payment.PaidAmount >= payment.Amount {
		payment.Status = "paid"
		payment.PaidAt = &now
	} else {
		payment.Status = "partial"
	}

	return s.repo.Update(payment)
}

// GetPaymentDetail gets payment detail with transactions
func (s *PaymentService) GetPaymentDetail(schoolCode string, paymentID uint) (*dto.PaymentDetailResponse, error) {
	payment, err := s.repo.FindByID(schoolCode, paymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	transactions, err := s.repo.GetTransactions(paymentID)
	if err != nil {
		return nil, err
	}

	var txResponses []dto.TransactionResponse
	for _, tx := range transactions {
		txResponses = append(txResponses, dto.ToTransactionResponse(&tx))
	}

	return &dto.PaymentDetailResponse{
		Payment:      dto.ToPaymentResponse(payment),
		Transactions: txResponses,
	}, nil
}
