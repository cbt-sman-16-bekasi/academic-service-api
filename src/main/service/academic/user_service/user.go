package user_service

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/user_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
)

type UserService struct {
	userRepo *school_repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: school_repository.NewUserRepository(),
	}
}

func (s *UserService) CreateNewUser(user *user.User) *user.User {
	isExist := s.userRepo.ReadUser(user.Username)
	if isExist != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", user.Username)))
	}

	if user.RoleUser.Code != "STUDENT" {
		password, salt, err := helper.HashPasswordArgon2(user.Password)
		if err != nil {
			panic(err)
		}
		user.Password = password
		user.Salt = salt
	}

	s.userRepo.Database.Create(&user)
	return user
}

func (s *UserService) UpdateUser(userId uint, user *user.User) *user.User {
	isExist := s.userRepo.FindById(userId)
	if isExist == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "User not found"))
	}

	if isExist.Username != user.Username {
		findUser := s.userRepo.ReadUser(user.Username)
		if findUser != nil {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", user.Username)))
		}
		isExist.Username = user.Username
	}

	if user.RoleUser.Code != "STUDENT" {
		password, salt, err := helper.HashPasswordArgon2(user.Password)
		if err != nil {
			panic(err)
		}
		isExist.Password = password
		isExist.Salt = salt
	}

	isExist.Status = user.Status
	s.userRepo.Database.Save(&isExist)
	return isExist
}

func (s *UserService) GetAllRole() []user.Role {
	return s.userRepo.AllRole()
}

func (s *UserService) GetAllUser(request pagination.Request[map[string]interface{}]) *database.Paginator {

	return database.NewPagination[map[string]interface{}]().
		SetModal([]view.VUser{}).
		SetRequest(&request).
		FindAllPaging()
}

func (s *UserService) GetUserById(userId uint) (user *view.VUser) {
	s.userRepo.Database.Where("id = ?", userId).First(&user)
	return
}

func (s *UserService) EnhanceSimpleUser(userId uint, request user_request.UserUpdateRequest) user_request.UserUpdateRequest {
	var userData *user.User
	s.userRepo.Database.Where("id = ?", userId).First(&userData)
	if userData == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "User not found"))
	}
	if userData.Username != request.Username {
		isExist := s.userRepo.ReadUser(request.Username)
		if isExist != nil {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", request.Username)))
		}
		userData.Username = request.Username
	}

	userData.Name = request.Name
	userData.Status = uint(request.Status)

	s.userRepo.Database.Save(&userData)
	return request
}

func (s *UserService) ResetPasswordUser(userId uint) {
	var userData = s.userRepo.FindById(userId)
	if userData == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed reset password, data user is not found or inactive"))
	}
	password, salt, err := helper.HashPasswordArgon2(userData.Username)
	if err != nil {
		panic(err)
	}
	userData.Password = password
	userData.Salt = salt

	s.userRepo.Database.Save(&userData)
}

func (s *UserService) CreateUser(request user_request.UserUpdateRequest) user_request.UserUpdateRequest {
	existUser := s.userRepo.ReadUser(request.Username)
	if existUser != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User already exists: %s", request.Username)))
	}

	s.CreateNewUser(&user.User{
		SchoolCode: "db74a42e-23a7-4cd2-bbe5-49cf79f86453",
		Username:   request.Username,
		Role:       uint(*request.Role),
		Status:     uint(request.Status),
		Name:       request.Name,
	})
	return request
}
