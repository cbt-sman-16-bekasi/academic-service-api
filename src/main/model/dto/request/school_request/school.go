package school_request

type ModifySchoolRequest struct {
	Logo              string `json:"logo"`
	LevelOfEducation  string `json:"level_of_education"`
	SchoolName        string `json:"school_name"`
	Nss               string `json:"nss"`
	Npsn              string `json:"npsn"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Address           string `json:"address"`
	Banner            string `json:"banner"`
	PrincipalName     string `json:"principal_name"`
	VicePrincipalName string `json:"vice_principal_name"`
	PrincipalNIP      string `json:"principal_nip"`
	VicePrincipalNIP  string `json:"vice_principal_nip"`
}
