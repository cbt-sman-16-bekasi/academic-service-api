package teacher_service

import (
	"fmt"
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
		SetPreloads("ClassSubject", "ClassSubject.Subject", "ClassSubject.Class").FindAllPaging()
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
		IsAccess: existTeacher.DetailUser.Status == 1,
		Gender:   existTeacher.Gender,
	}
}

func (t *TeacherService) CreateTeacher(request teacher_request.TeacherModifyRequest) teacher_response.TeacherDetailResponse {
	role := t.userRepository.ReadRole("TEACHER")
	if role == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Undefined role '%s'", "Teacher")))
	}

	var existingUser teacher.Teacher
	t.teacherRepository.Database.Where("nuptk = ?", request.Nuptk).First(&existingUser)
	if existingUser.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User with nuptk '%s' already exists", request.Nuptk)))
	}

	userService := user_service.NewUserService()
	useHasAccess := 1
	if request.IsAccess == false {
		useHasAccess = 0
	}
	resultUser := userService.CreateNewUser(&user.User{
		Username:   request.Nuptk,
		Role:       role.ID,
		Status:     uint(useHasAccess),
		Password:   request.Nuptk,
		SchoolCode: "db74a42e-23a7-4cd2-bbe5-49cf79f86453",
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
		Username: request.Nuptk,
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

	userService := user_service.NewUserService()
	useHasAccess := 1
	if request.IsAccess == false {
		useHasAccess = 0
	}
	_ = userService.UpdateUser(existTeacher.UserId, &user.User{
		Username: request.Nuptk,
		Status:   uint(useHasAccess),
		Password: request.Nuptk,
	})

	existTeacher.Name = request.Name
	existTeacher.Gender = request.Gender
	t.teacherRepository.Database.Save(&existTeacher)

	return teacher_response.TeacherDetailResponse{
		Nuptk:    request.Nuptk,
		Name:     request.Name,
		Username: request.Nuptk,
	}
}

func (t *TeacherService) DeleteById(id uint) {
	t.teacherRepository.Delete(id)
}

func (t *TeacherService) GetAllSubjectClassTeacher(teacherId uint) []teacher.TeacherClassSubject {
	var subjects []teacher.TeacherClassSubject
	t.teacherRepository.Database.Where("teacher_id = ?", teacherId).
		Preload("Subject").
		Preload("Class").
		Find(&subjects)

	return subjects
}

func (t *TeacherService) GetDetailSubjectClassTeacher(id uint) teacher.TeacherClassSubject {
	var subject teacher.TeacherClassSubject
	t.teacherRepository.Database.Where("id = ?", id).First(&subject)
	return subject
}

func (t *TeacherService) CreateSubjectClassTeacher(request teacher_request.TeacherMappingSubjectClass) teacher_request.TeacherMappingSubjectClass {
	t.validateCheckDuplicate(request)

	var multiple []teacher.TeacherClassSubject
	for _, u := range request.ClassId {
		data := teacher.TeacherClassSubject{
			TeacherId:   request.TeacherId,
			SubjectCode: request.SubjectId,
			ClassId:     u,
		}
		multiple = append(multiple, data)
	}

	t.teacherRepository.Database.Create(&multiple)

	return request
}

func (t *TeacherService) validateCheckDuplicate(request teacher_request.TeacherMappingSubjectClass) {
	for _, u := range request.ClassId {
		var checkExistRegister teacher.TeacherClassSubject
		t.teacherRepository.Database.Where("class_id = ? AND subject_code = ?", u, request.SubjectId).First(&checkExistRegister)

		if checkExistRegister.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Teacher with class id '%d' and '%s' already exists", request.ClassId, request.SubjectId)))
		}
	}
}

func (t *TeacherService) DeleteSubjectClassTeacher(id uint) {
	var subject teacher.TeacherClassSubject
	t.teacherRepository.Database.Where("id = ?", id).First(&subject)
	if subject.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Teacher with id '%d' not found", id)))
	}
	t.teacherRepository.Database.Delete(&subject)
}

func (t *TeacherService) UpdateSubjectClassTeacher(id uint, request teacher_request.TeacherMappingSubjectClass) teacher_request.TeacherMappingSubjectClass {
	//var subject teacher.TeacherClassSubject
	//t.teacherRepository.Database.Where("id = ?", id).First(&subject)
	//if subject.ID == 0 {
	//	panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Teacher with id '%d' not found", id)))
	//}
	//
	//if request.SubjectId != subject.SubjectCode {
	//	t.validateCheckDuplicate(request)
	//}
	//
	//if request.ClassId != subject.ClassId {
	//	t.validateCheckDuplicate(request)
	//}
	//
	//subject.SubjectCode = request.SubjectId
	//subject.ClassId = request.ClassId
	//
	//t.teacherRepository.Database.Save(&subject)
	return request
}
