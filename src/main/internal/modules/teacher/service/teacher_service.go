package service

import (
	"fmt"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/teacher/repository"
	userService "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/user/service"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/gin-gonic/gin"
)

type TeacherService struct {
	repo        *repository.TeacherRepository
	userService *userService.UserService
}

// NewTeacherService creates service with injected dependencies
func NewTeacherService(repo *repository.TeacherRepository, userSvc *userService.UserService) *TeacherService {
	return &TeacherService{repo: repo, userService: userSvc}
}

// GetAllTeacher returns paginated teachers
func (s *TeacherService) GetAllTeacher(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)

	// Filter by school_code for multi-tenant isolation
	filter := map[string]interface{}{}
	if request.Filter != nil {
		filter = *request.Filter
	}
	filter["school_code"] = claims.SchoolCode
	request.Filter = &filter

	return database.NewPagination[map[string]interface{}](s.repo.DB()).
		SetRequest(&request).
		SetModel([]teacher.Teacher{}).
		SetPreloads("ClassSubject", "ClassSubject.Subject", "ClassSubject.Class").
		FindAllPaging()
}

// DetailTeacher returns teacher detail
func (s *TeacherService) DetailTeacher(id uint) dto.TeacherDetailResponse {
	t, err := s.repo.FindByID(id)
	if err != nil || t.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Teacher not found"))
	}

	return dto.TeacherDetailResponse{
		Nuptk:    t.Nuptk,
		Name:     t.Name,
		Username: t.DetailUser.Username,
		Gender:   t.Gender,
		Role: dto.GeneralLabelKeyResponse{
			Key:   t.DetailUser.RoleUser.Code,
			Label: t.DetailUser.RoleUser.Name,
		},
		IsAccess: t.DetailUser.Status == 1,
	}
}

// CreateTeacher creates a new teacher
func (s *TeacherService) CreateTeacher(c *gin.Context, req dto.TeacherModifyRequest) dto.TeacherDetailResponse {
	claims := jwt.GetDataClaims(c)

	// Get TEACHER role
	role := s.userService.GetRole("TEACHER")
	if role == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Undefined role 'TEACHER'"))
	}

	// Check if NUPTK already exists
	existing, _ := s.repo.FindByNuptk(claims.SchoolCode, req.Nuptk)
	if existing != nil && existing.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("User with nuptk '%s' already exists", req.Nuptk)))
	}

	// Create user first
	status := uint(1)
	if !req.IsAccess {
		status = 0
	}

	resultUser := s.userService.CreateNewUser(&user.User{
		Username:   req.Nuptk,
		Role:       role.ID,
		Status:     status,
		Password:   req.Nuptk,
		SchoolCode: claims.SchoolCode,
	})

	// Create teacher
	newTeacher := teacher.Teacher{
		UserId:     resultUser.ID,
		Name:       req.Name,
		Nuptk:      req.Nuptk,
		Gender:     req.Gender,
		SchoolCode: claims.SchoolCode,
	}

	if err := s.repo.Create(&newTeacher); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return dto.TeacherDetailResponse{
		Nuptk:    req.Nuptk,
		Name:     req.Name,
		Username: req.Nuptk,
		Gender:   req.Gender,
		Role: dto.GeneralLabelKeyResponse{
			Key:   role.Code,
			Label: role.Name,
		},
		IsAccess: req.IsAccess,
	}
}

// UpdateTeacher updates an existing teacher
func (s *TeacherService) UpdateTeacher(c *gin.Context, id uint, req dto.TeacherModifyRequest) dto.TeacherDetailResponse {
	claims := jwt.GetDataClaims(c)
	schoolCode := claims.SchoolCode
	existingTeacher, err := s.repo.FindByID(id)
	if err != nil || existingTeacher.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, fmt.Sprintf("Teacher with id '%d' not found", id)))
	}

	// Check NUPTK uniqueness if changed
	if existingTeacher.Nuptk != req.Nuptk {
		duplicate, _ := s.repo.FindByNuptk(schoolCode, req.Nuptk)
		if duplicate != nil && duplicate.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Teacher with nuptk '%s' already exists", req.Nuptk)))
		}
		existingTeacher.Nuptk = req.Nuptk
	}

	// Update user
	status := uint(1)
	if !req.IsAccess {
		status = 0
	}
	s.userService.UpdateUserEntity(existingTeacher.UserId, &user.User{
		Username: req.Nuptk,
		Status:   status,
		Password: req.Nuptk,
	})

	existingTeacher.Name = req.Name
	existingTeacher.Gender = req.Gender

	if err := s.repo.Update(existingTeacher); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return dto.TeacherDetailResponse{
		Nuptk:    req.Nuptk,
		Name:     req.Name,
		Username: req.Nuptk,
		Gender:   req.Gender,
		IsAccess: req.IsAccess,
	}
}

// DeleteTeacher deletes a teacher
func (s *TeacherService) DeleteTeacher(id uint) {
	if err := s.repo.Delete(id); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
}

// GetAllSubjectClassTeacher returns all class subjects for a teacher
func (s *TeacherService) GetAllSubjectClassTeacher(teacherID uint) []teacher.TeacherClassSubject {
	subjects, _ := s.repo.FindClassSubjectsByTeacher(teacherID)
	return subjects
}

// GetDetailSubjectClassTeacher returns a specific class subject
func (s *TeacherService) GetDetailSubjectClassTeacher(id uint) *teacher.TeacherClassSubject {
	subject, _ := s.repo.FindClassSubjectByID(id)
	return subject
}

// CreateSubjectClassTeacher creates class subject assignments for a teacher
func (s *TeacherService) CreateSubjectClassTeacher(req dto.TeacherMappingSubjectClass) dto.TeacherMappingSubjectClass {
	// Validate duplicates
	for _, classID := range req.ClassId {
		existing, _ := s.repo.FindClassSubjectDuplicate(classID, req.SubjectId, req.TeacherId)
		if existing != nil && existing.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest,
				fmt.Sprintf("Teacher with class id '%d' and subject '%s' already exists", classID, req.SubjectId)))
		}
	}

	var subjects []teacher.TeacherClassSubject
	for _, classID := range req.ClassId {
		subjects = append(subjects, teacher.TeacherClassSubject{
			TeacherId:   req.TeacherId,
			SubjectCode: req.SubjectId,
			ClassId:     classID,
		})
	}

	if err := s.repo.CreateClassSubjects(subjects); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return req
}

// DeleteSubjectClassTeacher deletes a class subject assignment
func (s *TeacherService) DeleteSubjectClassTeacher(id uint) {
	subject, _ := s.repo.FindClassSubjectByID(id)
	if subject == nil || subject.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, fmt.Sprintf("Class subject with id '%d' not found", id)))
	}

	if err := s.repo.DeleteClassSubject(id); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
}

// UpdateSubjectClassTeacher updates a class subject assignment (stub for now)
func (s *TeacherService) UpdateSubjectClassTeacher(id uint, req dto.TeacherMappingSubjectClass) dto.TeacherMappingSubjectClass {
	// TODO: Implement update logic
	return req
}
