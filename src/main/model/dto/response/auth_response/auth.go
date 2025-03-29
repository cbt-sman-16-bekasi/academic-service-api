package auth_response

import "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"

type AuthResponse struct {
	Token string     `json:"token"`
	Exp   int64      `json:"exp"`
	User  *user.User `json:"user"`
}
