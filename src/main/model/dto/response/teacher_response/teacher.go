package teacher_response

import "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"

type TeacherDetailResponse struct {
	Nuptk    string                           `json:"nuptk"`
	Name     string                           `json:"name"`
	Username string                           `json:"username"`
	Role     response.GeneralLabelKeyResponse `json:"role"`
}
