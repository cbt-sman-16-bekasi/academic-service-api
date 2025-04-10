package school_response

type DetailSchool struct {
	Id               string `json:"id"`
	Logo             string `json:"logo"`
	LevelOfEducation string `json:"level_of_education"`
	SchoolName       string `json:"school_name"`
	Nss              string `json:"nss"`
	Npsn             string `json:"npsn"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Address          string `json:"address"`
	Banner           string `json:"banner"`
}

type DashboardResponse struct {
	TotalClass       int `json:"total_class"`
	TotalSubject     int `json:"total_subject"`
	TotalStudent     int `json:"total_student"`
	TotalExam        int `json:"total_exam"`
	TotalSessionExam int `json:"total_session_exam"`
	TotalReportExam  int `json:"total_report_exam"`
	TotalAccess      int `json:"total_access"`
}
