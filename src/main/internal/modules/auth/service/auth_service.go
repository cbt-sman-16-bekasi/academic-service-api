package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/auth/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	repo *repository.AuthRepository
}

// NewAuthService creates service with injected repository
func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Login authenticates admin/teacher users
func (s *AuthService) Login(username, password string) dto.AuthResponse {
	dataUser, err := s.repo.FindUserByUsername(username)
	if err != nil || dataUser.ID == 0 {
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
		SchoolCode: dataUser.SchoolCode,
	})

	if err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	detailUser := s.getDetailUser(dataUser)
	jwt.SaveDetailUser(username, detailUser, time.Duration(24)*time.Hour)

	return dto.AuthResponse{
		Token:  token,
		User:   dataUser,
		Exp:    exp,
		Detail: detailUser,
	}
}

// getDetailUser returns user detail based on role
func (s *AuthService) getDetailUser(dataUser *user.User) interface{} {
	switch dataUser.RoleUser.Code {
	case "STUDENT":
		return map[string]interface{}{
			"name": strings.ReplaceAll(dataUser.Username, "_", " "),
			"ID":   dataUser.ID,
		}
	case "TEACHER":
		teacherData, err := s.repo.FindTeacherByUserID(dataUser.ID)
		if err != nil || teacherData.ID == 0 {
			panic(exception.NewBadRequestExceptionStruct(response.Unauthorized, "You don't have access. Please contact your administrator"))
		}

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

// ChangePassword changes user password
func (s *AuthService) ChangePassword(c *gin.Context, req dto.ChangePasswordRequest) {
	dataClaims := jwt.GetDataClaims(c)
	dataUser, err := s.repo.FindUserByUsername(dataClaims.Username)
	if err != nil || dataUser.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "User not found"))
	}

	if !helper.VerifyPasswordArgon2(req.CurrentPassword, dataUser.Password, dataUser.Salt) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Current password is incorrect"))
	}

	newPassword, newSalt, _ := helper.HashPasswordArgon2(req.NewPassword)
	dataUser.Password = newPassword
	dataUser.Salt = newSalt

	if err := s.repo.UpdateUser(dataUser); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
}

// ChangeProfile changes user profile
func (s *AuthService) ChangeProfile(c *gin.Context, req dto.ChangeProfileRequest) {
	dataClaims := jwt.GetDataClaims(c)
	dataUser, err := s.repo.FindUserByUsername(dataClaims.Username)
	if err != nil || dataUser.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "User not found"))
	}

	// Check username uniqueness if changed
	if dataUser.Username != req.Username {
		existing, _ := s.repo.FindUserByUsername(req.Username)
		if existing != nil && existing.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Username already exists"))
		}
		dataUser.Username = req.Username
	}

	dataUser.Name = req.FullName
	dataUser.ProfileURL = req.ProfileURL

	if err := s.repo.UpdateUser(dataUser); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	// Update teacher name if teacher
	if dataClaims.Role == "TEACHER" {
		teacherData, _ := s.repo.FindTeacherByUserID(dataUser.ID)
		if teacherData != nil && teacherData.ID != 0 {
			teacherData.Name = req.FullName
			s.repo.UpdateTeacher(teacherData)
		}
	}
}
