package dto

// ----- Request DTOs -----

// ModifyClassRequest for create/update class
type ModifyClassRequest struct {
	ClassCode string `json:"class_code" binding:"required"`
	ClassName string `json:"class_name" binding:"required"`
}

// ModifyClassSubject for class subject relations
type ModifyClassSubject struct {
	ClassCode   string `json:"class_code" binding:"required"`
	SubjectCode string `json:"subject_code" binding:"required"`
}

// ModifyClassMemberRequest for adding/removing class members
type ModifyClassMemberRequest struct {
	ClassId   uint   `json:"class_id" binding:"required"`
	StudentId []uint `json:"student_id" binding:"required"`
}

// ----- Response DTOs -----

// GeneralLabelKeyResponse for key-label pairs
type GeneralLabelKeyResponse struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// ClassCodeResponse for class code
type ClassCodeResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// DetailClassResponse for class detail
type DetailClassResponse struct {
	ID        uint                    `json:"id"`
	ClassCode GeneralLabelKeyResponse `json:"class_code"`
	ClassName string                  `json:"class_name"`
}

// DetailClassSubjectResponse for class subject relation
type DetailClassSubjectResponse struct {
	ID        uint                    `json:"id"`
	ClassCode GeneralLabelKeyResponse `json:"class_code"`
	Subject   GeneralLabelKeyResponse `json:"subject"`
}

