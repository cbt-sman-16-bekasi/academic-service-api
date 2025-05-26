package curriculum_request

type SubjectRequest struct {
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	ClassCode []string `json:"class_code"`
}
