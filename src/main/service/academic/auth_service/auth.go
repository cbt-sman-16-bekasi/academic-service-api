package auth_service

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/auth_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/server/response"
	"time"
)

type AuthService struct {
	userRepository *school_repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepository: school_repository.NewUserRepository(),
	}
}

func (r *AuthService) Login(username string, password string) auth_response.AuthResponse {
	user := r.userRepository.ReadUser(username)
	if user == nil {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Username or Password is incorrect"))
	}

	if user.RoleUser.Code != "ADMIN" {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "You don't have access, You are not admin!"))
	}

	if !helper.VerifyPasswordArgon2(password, user.Password, user.Salt) {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Username or Password is incorrect"))
	}

	exp := time.Now().Add(time.Hour * 24).Unix()
	token, err := jwt.GenerateJWT(jwt.Claims{
		Username:   username,
		Role:       user.RoleUser.Code,
		Permission: []string{"create", "update", "delete", "read", "list"},
		SchoolCode: "db74a42e-23a7-4cd2-bbe5-49cf79f86453",
	})

	if err != nil {
		panic(exception.NewIntenalServerExceptionStruct(response.ServerError, err.Error()))
	}
	return auth_response.AuthResponse{
		Token: token,
		User:  user,
		Exp:   exp,
	}
}
