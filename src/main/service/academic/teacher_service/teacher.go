package teacher_service

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/teacher_request"
	response2 "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/teacher_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/service/academic/user_service"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
)

type TeacherService struct {
	teacherRepository *school_repository.TeacherRepository
	userRepository    *school_repository.UserRepository
}

func NewTeacherService() *TeacherService {
	return &TeacherService{
		teacherRepository: school_repository.NewTeacherRepository(),
		userRepository:    school_repository.NewUserRepository(),
	}
}

func (t *TeacherService) GetAllTeacher(request pagination.Request[map[string]interface{}]) *database.Paginator {
	paging := database.NewPagination[map[string]interface{}]().
		SetRequest(&request).
		SetModal([]teacher.Teacher{}).
		SetPreloads("DetailUser", "DetailUser.RoleUser").FindAllPaging()
	return paging
}

func (t *TeacherService) DetailTeacher(id uint) teacher_response.TeacherDetailResponse {
	existTeacher := t.teacherRepository.FindById(id)
	return teacher_response.TeacherDetailResponse{
		Nuptk:    existTeacher.Nuptk,
		Name:     existTeacher.Name,
		Username: existTeacher.DetailUser.Username,
		Role: response2.GeneralLabelKeyResponse{
			Key:   existTeacher.DetailUser.RoleUser.Code,
			Label: existTeacher.DetailUser.RoleUser.Name,
		},
	}
}

func (t *TeacherService) CreateTeacher(request teacher_request.TeacherModifyRequest) teacher_response.TeacherDetailResponse {
	role := t.userRepository.ReadRole(request.Role)
	if role == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Undefined role '%s'", request.Role)))
	}

	var existingUser teacher.Teacher
	t.teacherRepository.Database.Where("nuptk = ?", request.Nuptk).First(&existingUser)
	if existingUser.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User with nuptk '%s' already exists", request.Nuptk)))
	}

	userService := user_service.NewUserService()
	password, salt, _ := helper.HashPasswordArgon2(request.Password)
	resultUser := userService.CreateNewUser(&user.User{
		Username: request.Username,
		Role:     role.ID,
		Status:   1,
		Password: password,
		Salt:     salt,
	})

	teach := teacher.Teacher{
		UserId: resultUser.ID,
		Name:   request.Name,
		Nuptk:  request.Nuptk,
	}

	t.teacherRepository.Database.Save(&teach)

	return teacher_response.TeacherDetailResponse{
		Nuptk:    request.Nuptk,
		Name:     request.Name,
		Username: request.Username,
		Role: response2.GeneralLabelKeyResponse{
			Key:   role.Code,
			Label: role.Name,
		},
	}
}

func (t *TeacherService) UpdateTeacher(id uint, request teacher_request.TeacherModifyRequest) teacher_response.TeacherDetailResponse {
	existTeacher := t.teacherRepository.FindById(id)
	if existTeacher == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Teacher with id '%d' not found", id)))
	}

	if existTeacher.Nuptk != request.Nuptk {
		var teachNuptk teacher.Teacher
		t.teacherRepository.Database.Where("nuptk = ?", request.Nuptk).First(&teachNuptk)
		if teachNuptk.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Teacher with nuptk '%s' already exists", request.Nuptk)))
		}
		existTeacher.Nuptk = request.Nuptk
	}

	role := t.userRepository.ReadRole(request.Role)
	if role == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Undefined role '%s'", request.Role)))
	}

	userService := user_service.NewUserService()
	password, salt, _ := helper.HashPasswordArgon2(request.Password)
	_ = userService.UpdateUser(existTeacher.UserId, &user.User{
		Username: request.Username,
		Role:     role.ID,
		Status:   1,
		Password: password,
		Salt:     salt,
	})

	t.teacherRepository.Database.Save(&existTeacher)

	return teacher_response.TeacherDetailResponse{
		Nuptk:    request.Nuptk,
		Name:     request.Name,
		Username: request.Username,
		Role: response2.GeneralLabelKeyResponse{
			Key:   role.Code,
			Label: role.Name,
		},
	}
}

func (t *TeacherService) DeleteById(id uint) {
	t.teacherRepository.Delete(id)
}
