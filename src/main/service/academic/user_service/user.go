package user_service

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/exception"
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
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "User already exists"))
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

	s.userRepo.Database.Save(&isExist)
	return isExist
}

func (s *UserService) GetAllRole() []user.Role {
	return s.userRepo.AllRole()
}
