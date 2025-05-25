package view

type DashboardSummary struct {
	TotalUsers         int `json:"total_users" gorm:"column:total_users"`
	TotalClasses       int `json:"total_classes" gorm:"column:total_classes"`
	TotalClassSubjects int `json:"total_class_subjects" gorm:"column:total_class_subjects"`
	TotalStudents      int `json:"total_students" gorm:"column:total_students"`
	TotalExams         int `json:"total_exams" gorm:"column:total_exams"`
	TotalExamSessions  int `json:"total_exam_sessions" gorm:"column:total_exam_sessions"`
	TotalReport        int `json:"total_report" gorm:"column:total_report"`
}

func (d *DashboardSummary) TableName() string {
	return "public.v_dashboard_summary"
}

type DashboardTeacher struct {
	TotalUsers         int `json:"total_users" gorm:"column:total_users"`
	TotalClasses       int `json:"total_classes" gorm:"column:total_classes"`
	TotalClassSubjects int `json:"total_class_subjects" gorm:"column:total_class_subjects"`
	TotalStudents      int `json:"total_students" gorm:"column:total_students"`
	TotalExams         int `json:"total_exams" gorm:"column:total_exams"`
	TotalExamSessions  int `json:"total_exam_sessions" gorm:"column:total_exam_sessions"`
	TotalReport        int `json:"total_report" gorm:"column:total_report"`
}

func (d *DashboardTeacher) TableName() string {
	return "public.v_dashboard_teacher"
}
