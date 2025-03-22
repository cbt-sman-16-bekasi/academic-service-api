package student_response

import "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"

type DetailStudentResponse struct {
	Nisn   string                           `json:"nisn"`
	Name   string                           `json:"name"`
	Gender string                           `json:"gender"`
	Class  response.GeneralLabelKeyResponse `json:"class"`
}
