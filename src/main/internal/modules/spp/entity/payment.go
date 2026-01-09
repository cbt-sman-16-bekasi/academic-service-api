package entity

import (
	"time"

	"gorm.io/gorm"
)

// Payment represents a student payment record
type Payment struct {
	gorm.Model
	SchoolCode   string     `json:"school_code" gorm:"index;not null"`
	StudentID    uint       `json:"student_id" gorm:"index;not null"`
	PaymentType  string     `json:"payment_type" gorm:"not null"` // SPP, DAFTAR_ULANG, etc
	Amount       float64    `json:"amount" gorm:"not null"`
	PaidAmount   float64    `json:"paid_amount" gorm:"default:0"`
	Status       string     `json:"status" gorm:"default:'pending'"` // pending, partial, paid
	DueDate      time.Time  `json:"due_date"`
	PaidAt       *time.Time `json:"paid_at"`
	AcademicYear string     `json:"academic_year"`
	Semester     int        `json:"semester"`
	Notes        string     `json:"notes"`
}

// PaymentTransaction represents payment transaction history
type PaymentTransaction struct {
	gorm.Model
	SchoolCode    string    `json:"school_code" gorm:"index;not null"`
	PaymentID     uint      `json:"payment_id" gorm:"index;not null"`
	Amount        float64   `json:"amount" gorm:"not null"`
	PaymentMethod string    `json:"payment_method"` // cash, transfer, qris
	TransactionID string    `json:"transaction_id"`
	PaidAt        time.Time `json:"paid_at"`
	CreatedBy     uint      `json:"created_by"`
	Notes         string    `json:"notes"`

	// Relations
	Payment Payment `json:"payment" gorm:"foreignKey:PaymentID"`
}

// PaymentConfig represents SPP configuration per school
type PaymentConfig struct {
	gorm.Model
	SchoolCode    string  `json:"school_code" gorm:"uniqueIndex;not null"`
	DefaultAmount float64 `json:"default_amount"`
	DueDateDay    int     `json:"due_date_day" gorm:"default:10"` // tanggal jatuh tempo
	GracePeriod   int     `json:"grace_period" gorm:"default:7"`  // hari
	LateFee       float64 `json:"late_fee" gorm:"default:0"`
}

func (Payment) TableName() string {
	return "spp_payments"
}

func (PaymentTransaction) TableName() string {
	return "spp_transactions"
}

func (PaymentConfig) TableName() string {
	return "spp_configs"
}
