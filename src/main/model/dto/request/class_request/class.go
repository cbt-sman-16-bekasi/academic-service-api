package class_request

type ModifyClassRequest struct {
	ClassCode string `json:"class_code"`
	ClassName string `json:"class_name"`
}

type ModifyClassSubject struct {
	ClassCode   string `json:"class_code"`
	SubjectCode string `json:"subject_code"`
}
