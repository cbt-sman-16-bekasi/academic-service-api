package class_response

import "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"

type ClassCodeResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type DetailClassResponse struct {
	ID        uint                             `json:"id"`
	ClassCode response.GeneralLabelKeyResponse `json:"class_code"`
	ClassName string                           `json:"class_name"`
}
