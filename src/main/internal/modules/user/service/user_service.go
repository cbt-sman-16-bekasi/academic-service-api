package service

import (
	"fmt"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates service with injected repository
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetAllRole returns all roles
func (s *UserService) GetAllRole() []user.Role {
	roles, _ := s.repo.AllRoles()
	return roles
}

// GetAllUser returns paginated users
func (s *UserService) GetAllUser(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)

	// Filter by school_code for multi-tenant isolation
	filter := map[string]interface{}{}
	if request.Filter != nil {
		filter = *request.Filter
	}
	filter["school_code"] = claims.SchoolCode
	request.Filter = &filter

	return database.NewPagination[map[string]interface{}](s.repo.DB()).
		SetModel([]view.VUser{}).
		SetRequest(&request).
		FindAllPaging()
}

// GetUserById returns user by ID
func (s *UserService) GetUserById(userId uint) *view.VUser {
	u, err := s.repo.FindViewByID(userId)
	if err != nil || u.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "User not found"))
	}
	return u
}

// EnhanceSimpleUser updates user data
func (s *UserService) EnhanceSimpleUser(userId uint, req dto.UserUpdateRequest) dto.UserUpdateRequest {
	userData, err := s.repo.FindByID(userId)
	if err != nil || userData.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "User not found"))
	}

	// Check username uniqueness if changed
	if userData.Username != req.Username {
		existing, _ := s.repo.FindByUsername(req.Username)
		if existing != nil && existing.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", req.Username)))
		}
		userData.Username = req.Username
	}

	userData.Name = req.Name
	userData.Status = uint(req.Status)

	if err := s.repo.Update(userData); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return req
}

// ResetPasswordUser resets user password to username
func (s *UserService) ResetPasswordUser(userId uint) {
	userData, err := s.repo.FindByID(userId)
	if err != nil || userData.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "User not found"))
	}

	password, salt, err := helper.HashPasswordArgon2(userData.Username)
	if err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
	userData.Password = password
	userData.Salt = salt

	if err := s.repo.Update(userData); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(c *gin.Context, req dto.UserUpdateRequest) dto.UserUpdateRequest {
	claims := jwt.GetDataClaims(c)

	// Check if username exists
	existing, _ := s.repo.FindByUsername(req.Username)
	if existing != nil && existing.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", req.Username)))
	}

	// Hash password
	password, salt, err := helper.HashPasswordArgon2(req.Username)
	if err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	newUser := &user.User{
		SchoolCode: claims.SchoolCode,
		Username:   req.Username,
		Name:       req.Name,
		Role:       uint(*req.Role),
		Status:     uint(req.Status),
		Password:   password,
		Salt:       salt,
	}

	if err := s.repo.Create(newUser); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return req
}

// ========================================
// Cross-module helper methods
// Used by other modules (student, teacher)
// ========================================

// CreateNewUser creates a new user entity (for cross-module use)
func (s *UserService) CreateNewUser(u *user.User) *user.User {
	// Check if username exists
	existing, _ := s.repo.FindByUsername(u.Username)
	if existing != nil && existing.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", u.Username)))
	}

	// Hash password if role is not STUDENT
	role, _ := s.repo.FindRoleByCode(u.RoleUser.Code)
	if role == nil || role.Code != "STUDENT" {
		if u.Password != "" {
			password, salt, err := helper.HashPasswordArgon2(u.Password)
			if err != nil {
				panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
			}
			u.Password = password
			u.Salt = salt
		}
	}

	if err := s.repo.Create(u); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return u
}

// UpdateUserEntity updates user entity (for cross-module use)
func (s *UserService) UpdateUserEntity(userId uint, u *user.User) *user.User {
	existing, err := s.repo.FindByID(userId)
	if err != nil || existing.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "User not found"))
	}

	// Check username uniqueness if changed
	if existing.Username != u.Username {
		duplicate, _ := s.repo.FindByUsername(u.Username)
		if duplicate != nil && duplicate.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", u.Username)))
		}
		existing.Username = u.Username
	}

	// Hash password if provided and role is not STUDENT
	if u.Password != "" {
		role, _ := s.repo.FindRoleByCode(existing.RoleUser.Code)
		if role == nil || role.Code != "STUDENT" {
			password, salt, err := helper.HashPasswordArgon2(u.Password)
			if err != nil {
				panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
			}
			existing.Password = password
			existing.Salt = salt
		}
	}

	existing.Status = u.Status

	if err := s.repo.Update(existing); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return existing
}

// GetRole returns role by code (for cross-module use)
func (s *UserService) GetRole(code string) *user.Role {
	role, err := s.repo.FindRoleByCode(code)
	if err != nil || role.ID == 0 {
		return nil
	}
	return role
}
