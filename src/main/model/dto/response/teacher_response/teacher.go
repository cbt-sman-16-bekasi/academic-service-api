package teacher_response

import "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"

type TeacherDetailResponse struct {
	Nuptk    string                           `json:"nuptk"`
	Name     string                           `json:"name"`
	Username string                           `json:"username"`
	Gender   string                           `json:"gender"`
	Role     response.GeneralLabelKeyResponse `json:"role"`
	IsAccess bool                             `json:"isAccess"`
}
