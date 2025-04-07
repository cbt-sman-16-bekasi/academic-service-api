package auth_request

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CBTAuthRequest struct {
	Username string `json:"username"`
}
