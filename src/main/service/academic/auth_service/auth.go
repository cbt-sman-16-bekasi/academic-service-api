package auth_service

import (
	"encoding/json"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/auth_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/auth_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/server/response"
	"strings"
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
	dataUser := r.userRepository.ReadUser(username)
	if dataUser == nil {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Username or Password is incorrect"))
	}

	if dataUser.RoleUser.Code == "STUDENT" {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "You don't have access, You are not admin!"))
	}

	if !helper.VerifyPasswordArgon2(password, dataUser.Password, dataUser.Salt) {
		panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "Username or Password is incorrect"))
	}

	exp := time.Now().Add(time.Hour * 24).Unix()
	token, err := jwt.GenerateJWT(jwt.Claims{
		Username:   username,
		Role:       dataUser.RoleUser.Code,
		Permission: []string{"create", "update", "delete", "read", "list"},
		SchoolCode: "db74a42e-23a7-4cd2-bbe5-49cf79f86453",
	})

	if err != nil {
		panic(exception.NewIntenalServerExceptionStruct(response.ServerError, err.Error()))
	}
	detailUser := r.detailUser(dataUser)
	jwt.SaveDetailUser(username, detailUser, time.Duration(24)*time.Hour)
	return auth_response.AuthResponse{
		Token:  token,
		User:   dataUser,
		Exp:    exp,
		Detail: detailUser,
	}
}

func (r *AuthService) detailUser(dataUser *user.User) interface{} {
	switch dataUser.RoleUser.Code {
	case "STUDENT":
		return map[string]interface{}{
			"name": strings.ReplaceAll(dataUser.Username, "_", " "),
			"ID":   dataUser.ID,
		}
	case "TEACHER":
		var teacherData teacher.Teacher
		r.userRepository.Database.Where("user_id", dataUser.ID).First(&teacherData)
		if teacherData.ID == 0 {
			panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "You don't have access. Please contact your administrator"))
		}

		// Konversi ke map[string]interface{}
		var result map[string]interface{}
		data, _ := json.Marshal(teacherData)
		_ = json.Unmarshal(data, &result)

		result["profile_url"] = dataUser.ProfileURL
		return result
	case "ADMIN":
		return map[string]interface{}{
			"name":        strings.ReplaceAll(dataUser.Username, "_", " "),
			"ID":          dataUser.ID,
			"profile_url": dataUser.ProfileURL,
		}
	default:
		return nil
	}
}

func (r *AuthService) ChangePassword(c *gin.Context, request auth_request.ChangePasswordRequest) {
	dataClaims := jwt.GetDataClaims(c)
	dataUser := r.userRepository.ReadUser(dataClaims.Username)

	if !helper.VerifyPasswordArgon2(request.CurrentPassword, dataUser.Password, dataUser.Salt) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Current password is incorrect"))
	}

	newPassword, newSalt, _ := helper.HashPasswordArgon2(request.NewPassword)
	dataUser.Password = newPassword
	dataUser.Salt = newSalt

	r.userRepository.Database.Save(dataUser)
}

func (r *AuthService) ChangeProfile(c *gin.Context, request auth_request.ChangeProfileRequest) {
	dataClaims := jwt.GetDataClaims(c)
	dataUser := r.userRepository.ReadUser(dataClaims.Username)
	if dataUser.Username != request.Username {
		existUsername := r.userRepository.ReadUser(request.Username)
		if existUsername != nil {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Username already exists"))
		}
		dataUser.Username = request.Username
	}
	dataUser.Name = request.FullName
	dataUser.ProfileURL = request.ProfileURL

	r.userRepository.Database.Save(dataUser)

	if dataClaims.Role == "TEACHER" {
		var teacherData teacher.Teacher
		r.userRepository.Database.Where("user_id", dataUser.ID).First(&teacherData)
		if teacherData.ID != 0 {
			teacherData.Name = request.FullName
			r.userRepository.Database.Save(teacherData)
		}
	}
}
