package school_response

type DetailSchool struct {
	Id         string `json:"id"`
	SchoolName string `json:"school_name"`
	Logo       string `json:"logo"`
	Address    string `json:"address"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type DashboardResponse struct {
	TotalClass       int `json:"total_class"`
	TotalSubject     int `json:"total_subject"`
	TotalStudent     int `json:"total_student"`
	TotalExam        int `json:"total_exam"`
	TotalSessionExam int `json:"total_session_exam"`
	TotalReportExam  int `json:"total_report_exam"`
}
