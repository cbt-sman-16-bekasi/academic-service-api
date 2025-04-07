package exam_request

type ModifyTypeExamRequest struct {
	TypeExam     string `json:"type_exam"`
	CodeTypeExam string `json:"code_type_exam"`
	Role         string `json:"role"`
	Color        string `json:"color"`
}
