package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/entity"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates repository with injected DB
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// FindByID finds payment by ID with school scope
func (r *PaymentRepository) FindByID(schoolCode string, id uint) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.Where("school_code = ? AND id = ?", schoolCode, id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByStudent finds all payments for a student
func (r *PaymentRepository) FindByStudent(schoolCode string, studentID uint) ([]entity.Payment, error) {
	var payments []entity.Payment
	err := r.db.Where("school_code = ? AND student_id = ?", schoolCode, studentID).
		Order("due_date DESC").
		Find(&payments).Error
	return payments, err
}

// FindPending finds all pending payments for a school
func (r *PaymentRepository) FindPending(schoolCode string) ([]entity.Payment, error) {
	var payments []entity.Payment
	err := r.db.Where("school_code = ? AND status IN ?", schoolCode, []string{"pending", "partial"}).
		Order("due_date ASC").
		Find(&payments).Error
	return payments, err
}

// Create creates a new payment
func (r *PaymentRepository) Create(payment *entity.Payment) error {
	return r.db.Create(payment).Error
}

// Update updates a payment
func (r *PaymentRepository) Update(payment *entity.Payment) error {
	return r.db.Save(payment).Error
}

// Delete soft deletes a payment
func (r *PaymentRepository) Delete(schoolCode string, id uint) error {
	return r.db.Where("school_code = ? AND id = ?", schoolCode, id).Delete(&entity.Payment{}).Error
}

// CreateTransaction creates a payment transaction
func (r *PaymentRepository) CreateTransaction(tx *entity.PaymentTransaction) error {
	return r.db.Create(tx).Error
}

// GetTransactions gets all transactions for a payment
func (r *PaymentRepository) GetTransactions(paymentID uint) ([]entity.PaymentTransaction, error) {
	var transactions []entity.PaymentTransaction
	err := r.db.Where("payment_id = ?", paymentID).Order("paid_at DESC").Find(&transactions).Error
	return transactions, err
}
