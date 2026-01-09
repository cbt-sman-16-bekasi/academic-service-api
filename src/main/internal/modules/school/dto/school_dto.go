package dto

// ----- Request DTOs -----

// ModifySchoolRequest for updating school information
type ModifySchoolRequest struct {
	SchoolName        string `json:"school_name"`
	Logo              string `json:"logo"`
	Nss               string `json:"nss"`
	Npsn              string `json:"npsn"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Address           string `json:"address"`
	Banner            string `json:"banner"`
	PrincipalName     string `json:"principal_name"`
	VicePrincipalName string `json:"vice_principal_name"`
	PrincipalNIP      string `json:"principal_nip"`
	VicePrincipalNIP  string `json:"vice_principal_nip"`
}

// ModifyClassSubjectRequest for creating/updating class-subject mapping
type ModifyClassSubjectRequest struct {
	ClassCode   string `json:"class_code" binding:"required"`
	SubjectCode string `json:"subject_code" binding:"required"`
}

// ----- Response DTOs -----

// DetailSchoolResponse for school information
type DetailSchoolResponse struct {
	Id                string `json:"id"`
	SchoolName        string `json:"school_name"`
	Address           string `json:"address"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Logo              string `json:"logo"`
	Npsn              string `json:"npsn"`
	Nss               string `json:"nss"`
	Banner            string `json:"banner"`
	PrincipalName     string `json:"principal_name"`
	PrincipalNIP      string `json:"principal_nip"`
	VicePrincipalName string `json:"vice_principal_name"`
	VicePrincipalNIP  string `json:"vice_principal_nip"`
	LevelOfEducation  string `json:"level_of_education"`
}

// ClassCodeResponse for class code listing
type ClassCodeResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// SubjectResponse for subject listing
type SubjectResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Subject     string `json:"subject"`
	SubjectType string `json:"subject_type"`
	Description string `json:"description"`
}

// GeneralLabelKeyResponse for key-label pairs
type GeneralLabelKeyResponse struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// DetailClassSubjectResponse for class-subject detail
type DetailClassSubjectResponse struct {
	ID        uint                    `json:"id"`
	ClassCode GeneralLabelKeyResponse `json:"class_code"`
	Subject   GeneralLabelKeyResponse `json:"subject"`
}

// DashboardResponse for dashboard data
type DashboardResponse struct {
	TotalClass       int `json:"total_class"`
	TotalSubject     int `json:"total_subject"`
	TotalStudent     int `json:"total_student"`
	TotalExam        int `json:"total_exam"`
	TotalSessionExam int `json:"total_session_exam"`
	TotalReportExam  int `json:"total_report_exam"`
	TotalAccess      int `json:"total_access"`
}

// SchoolBriefInfo for config response
type SchoolBriefInfo struct {
	SchoolCode string `json:"school_code"`
	SchoolName string `json:"school_name"`
	Logo       string `json:"logo"`
	Banner     string `json:"banner"`
	Address    string `json:"address"`
}

// ConfigurationResponse for system config
type ConfigurationResponse struct {
	ClientID    string           `json:"client_id"`
	ThemePreset string           `json:"theme_preset,omitempty"`
	ThemeCustom string           `json:"theme_custom,omitempty"`
	SchoolInfo  *SchoolBriefInfo `json:"school_info,omitempty"`
}
