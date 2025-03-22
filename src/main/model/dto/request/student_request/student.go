package student_request

type StudentModifyRequest struct {
	Nisn    string `json:"nisn"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	ClassId uint   `json:"class_id"`
}
