package dto

// SubjectRequest for create/update subject
type SubjectRequest struct {
	Name      string   `json:"name" binding:"required"`
	Code      string   `json:"code" binding:"required"`
	ClassCode []string `json:"class_code"`
}

// SubjectResponse for subject response
type SubjectResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Subject     string `json:"subject"`
	SubjectType string `json:"subject_type,omitempty"`
	Description string `json:"description,omitempty"`
	ClassCode   string `json:"class_code,omitempty"`
	SchoolCode  string `json:"school_code,omitempty"`
}

