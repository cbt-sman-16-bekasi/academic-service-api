package cbt

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ExamSessionReport struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	SessionID    string    `json:"session_id" gorm:"not null"`
	ExamCode     string    `json:"exam_code" gorm:"type:varchar(50)"`
	SessionName  string    `json:"session_name"`
	ExamName     string    `json:"exam_name"`
	Subject      string    `json:"subject"`
	Kelas        string    `json:"kelas"`
	Total        int       `json:"total"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Status       string    `json:"status" gorm:"type:text;check:status IN ('AKTIF','SELESAI')"`
	CreatedBy    uint      `json:"created_by"`
	StatusReport string    `json:"status_report" gorm:"type:varchar(50)"`
	ReportURL    string    `json:"report_url"`
	gorm.Model
}

func (e *ExamSessionReport) TableName() string {
	return "cbt_service.exam_session_report"
}

func (e *ExamSessionReport) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	return nil
}
