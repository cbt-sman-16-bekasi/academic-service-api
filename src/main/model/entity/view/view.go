package view

type MasterBankQuestionResponse struct {
	ID           uint   `json:"id"`
	Subject      string `json:"subject"`
	ClassCode    string `json:"class_code"`
	Total        int    `json:"total"`
	TypeQuestion string `json:"type_question"`
	Name         string `json:"name"`
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
