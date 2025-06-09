package view

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"time"
)

type MasterBankQuestionResponse struct {
	ID           uint   `json:"id"`
	Subject      string `json:"subject"`
	SubjectCode  string `json:"subject_code"`
	ClassCode    string `json:"class_code"`
	Total        int    `json:"total"`
	TypeQuestion string `json:"type_question"`
	Name         string `json:"name"`
	BankName     string `json:"bank_name"`
}

func (s *MasterBankQuestionResponse) TableName() string {
	return "public.v_master_bank_question"
}

type SummaryExamSession struct {
	ID                 uint   `gorm:"column:id"`
	SessionID          string `gorm:"column:session_id"`
	TotalLogin         int64  `gorm:"column:total_login"`
	TotalStudentSubmit int64  `gorm:"column:total_student_submit"`
	TotalCheating      int64  `gorm:"column:total_cheating"`
	TotalTimeIsOver    int64  `gorm:"column:total_time_is_over"`
}

func (s *SummaryExamSession) TableName() string {
	return "public.v_summary_exam_session"
}

type ExamSessionActiveToday struct {
	ID          uint      `json:"id"`
	SessionID   string    `json:"session_id"`
	Class       uint      `json:"class"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	SessionName string    `json:"session_name"`
	ExamName    string    `json:"exam_name"`
	SubjectName string    `json:"subject_name"`
	TypeExam    string    `json:"type_exam"`
}

func (e *ExamSessionActiveToday) TableName() string {
	return "public.v_exam_session_active"
}

type VStudent struct {
	ID        uint   `json:"id" gorm:"column:id"`
	UserID    uint   `json:"user_id" gorm:"column:user_id"`
	NISN      string `json:"nisn" gorm:"column:nisn"`
	Name      string `json:"name" gorm:"column:name"`
	Gender    string `json:"gender" gorm:"column:gender"`
	ClassName string `json:"class_name" gorm:"column:class_name"`
	ClassID   uint   `json:"class_id" gorm:"column:class_id"`
}

// TableName sets the name of the table/view for GORM
func (v *VStudent) TableName() string {
	return "public.v_student"
}

type ExamSessionView struct {
	ID           string    `gorm:"column:id" json:"id"`
	SessionID    string    `gorm:"column:session_id" json:"session_id"`
	SessionName  string    `gorm:"column:session_name" json:"session_name"`
	ExamName     string    `gorm:"column:exam_name" json:"exam_name"`
	Subject      *string   `gorm:"column:subject" json:"subject"` // nullable
	Kelas        *string   `gorm:"column:kelas" json:"kelas"`     // nullable (hasil string_agg)
	Total        *int      `gorm:"column:total" json:"total"`     // nullable (karena join bisa kosong)
	StartDate    time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate      time.Time `gorm:"column:end_date" json:"end_date"`
	Status       string    `gorm:"column:status" json:"status"`
	StatusReport string    `gorm:"column:status_report" json:"status_report"`
	ReportUrl    string    `gorm:"column:report_url" json:"report_url"`
	CreatedBy    uint      `gorm:"column:created_by" json:"created_by"`
}

func (e *ExamSessionView) TableName() string {
	return "public.v_exam_session"
}

type ExamSessionReadyReport struct {
	ID          string    `gorm:"column:id" json:"id"`
	SessionID   string    `gorm:"column:session_id" json:"session_id"`
	SessionName string    `gorm:"column:session_name" json:"session_name"`
	ExamName    string    `gorm:"column:exam_name" json:"exam_name"`
	TypeExam    string    `gorm:"column:type_exam" json:"type_exam"`
	Subject     string    `gorm:"column:subject" json:"subject"` // nullable
	Kelas       string    `gorm:"column:kelas" json:"kelas"`     // nullable (hasil string_agg)
	KelasID     string    `gorm:"column:kelasid" json:"kelasId"` // nullable (hasil string_agg)
	StartDate   time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate     time.Time `gorm:"column:end_date" json:"end_date"`
}

func (e *ExamSessionReadyReport) TableName() string {
	return "public.v_exam_session_ready_report"
}

type ExamSessionReportScoreView struct {
	ID           string    `gorm:"column:id" json:"id"`
	SessionID    string    `gorm:"column:session_id" json:"session_id"`
	SessionName  string    `gorm:"column:session_name" json:"session_name"`
	ExamName     string    `gorm:"column:exam_name" json:"exam_name"`
	Subject      string    `gorm:"column:subject" json:"subject"` // nullable
	Kelas        string    `gorm:"column:kelas" json:"kelas"`     // nullable (hasil string_agg)
	Total        int       `gorm:"column:total" json:"total"`     // nullable (karena join bisa kosong)
	StartDate    time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate      time.Time `gorm:"column:end_date" json:"end_date"`
	Status       string    `gorm:"column:status" json:"status"`
	CreatedBy    uint      `gorm:"column:created_by" json:"created_by"`
	ReportUrl    string    `gorm:"column:report_url" json:"report_url"`
	StatusReport string    `gorm:"column:status_report" json:"status_report"`
}

func (e *ExamSessionReportScoreView) TableName() string {
	return "public.v_exam_session_report_score"
}

type VClass struct {
	ID           uint   `json:"ID" gorm:"column:id"`
	ClassCode    string `json:"classCode" gorm:"column:class_code"`
	ClassName    string `json:"className" gorm:"column:class_name"`
	TotalStudent int    `json:"total_student"`
}

func (v *VClass) TableName() string {
	return "public.v_class"
}

type VUser struct {
	ID         int64     `json:"id" gorm:"column:id"`
	SchoolCode string    `json:"school_code" gorm:"column:school_code"`
	Username   string    `json:"username" gorm:"column:username"`
	Name       string    `json:"name" gorm:"column:name"`
	ProfileURL string    `json:"profile_url" gorm:"column:profile_url"`
	Status     int       `json:"status" gorm:"column:status"`
	RoleName   string    `json:"role_name" gorm:"column:role_name"`
	Role       int       `json:"role" gorm:"column:role"`
	LastUpdate time.Time `json:"last_update" gorm:"column:last_update"`
}

func (v *VUser) TableName() string {
	return "public.v_user"
}

type VSubject struct {
	curriculum.Subject
	ClassCode string `json:"class_code" gorm:"column:class_code"`
}

func (v *VSubject) TableName() string {
	return "public.v_subject"
}
