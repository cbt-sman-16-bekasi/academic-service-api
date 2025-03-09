package school_response

type DetailSchool struct {
	Id         string `json:"id"`
	SchoolName string `json:"school_name"`
	Logo       string `json:"logo"`
	Address    string `json:"address"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}
