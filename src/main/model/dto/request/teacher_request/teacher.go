package teacher_request

type TeacherModifyRequest struct {
	Nuptk    string `json:"nuptk"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
