package auth_request

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CBTAuthRequest struct {
	Username string `json:"username"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangeProfileRequest struct {
	FullName   string `json:"full_name"`
	Username   string `json:"username"`
	ProfileURL string `json:"profile_url"`
}
