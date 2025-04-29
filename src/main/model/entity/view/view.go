package view

type MasterBankQuestionResponse struct {
	ID           uint   `json:"id"`
	Subject      string `json:"subject"`
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

func (SummaryExamSession) TableName() string {
	return "public.v_summary_exam_session"
}

type ExamSessionActiveToday struct {
	ID        uint   `json:"id"`
	SessionID string `json:"session_id"`
	Class     uint   `json:"class"`
}

func (ExamSessionActiveToday) TableName() string {
	return "public.v_exam_session_active_today"
}

type VStudent struct {
	ID        uint   `json:"id" db:"id" gorm:"column:id"`
	UserID    uint   `json:"user_id" db:"user_id" gorm:"column:user_id"`
	NISN      string `json:"nisn" db:"nisn" gorm:"column:nisn"`
	Name      string `json:"name" db:"name" gorm:"column:name"`
	Gender    string `json:"gender" db:"gender" gorm:"column:gender"`
	ClassName string `json:"class_name" db:"class_name" gorm:"column:class_name"`
	ClassID   uint   `json:"class_id" db:"class_id" gorm:"column:class_id"`
}

// TableName sets the name of the table/view for GORM
func (VStudent) TableName() string {
	return "public.v_student"
}
