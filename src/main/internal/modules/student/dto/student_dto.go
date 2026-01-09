package dto

// ----- Request DTOs -----

// StudentModifyRequest for create/update student
type StudentModifyRequest struct {
	Nisn    string `json:"nisn" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Gender  string `json:"gender"`
	ClassId uint   `json:"class_id" binding:"required"`
}

// ----- Response DTOs -----

// GeneralLabelKeyResponse for key-label pairs
type GeneralLabelKeyResponse struct {
	Key   interface{} `json:"key"`
	Label string      `json:"label"`
}

// DetailStudentResponse for student detail
type DetailStudentResponse struct {
	Nisn   string                  `json:"nisn"`
	Name   string                  `json:"name"`
	Gender string                  `json:"gender"`
	Class  GeneralLabelKeyResponse `json:"class"`
}
