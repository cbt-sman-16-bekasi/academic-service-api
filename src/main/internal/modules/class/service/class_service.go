package service

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/class/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/gin-gonic/gin"
)

type ClassService struct {
	repo *repository.ClassRepository
}

// NewClassService creates service with injected repository
func NewClassService(repo *repository.ClassRepository) *ClassService {
	return &ClassService{repo: repo}
}

// FindAllClass returns paginated classes
func (s *ClassService) FindAllClass(ctx *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(ctx)

	// Filter by school_code for multi-tenant isolation
	filter := map[string]interface{}{}
	if request.Filter != nil {
		filter = *request.Filter
	}
	filter["school_code"] = claims.SchoolCode

	// If not admin, filter by teacher's classes
	if claims.Role != "ADMIN" {
		teacherClassSubjects, _ := s.repo.FindTeacherClassSubjects(uint(jwt.GetID(claims.Username)))

		var classIDs []interface{}
		for _, tcs := range teacherClassSubjects {
			classIDs = append(classIDs, tcs.ClassId)
		}
		filter["id"] = classIDs
	}

	request.Filter = &filter
	return database.NewPagination[map[string]interface{}](s.repo.DB()).
		SetModel([]view.VClass{}).
		SetRequest(&request).
		FindAllPaging()
}

// GetDetailClass returns class detail
func (s *ClassService) GetDetailClass(ctx *gin.Context, id uint) dto.DetailClassResponse {
	claims := jwt.GetDataClaims(ctx)

	class, err := s.repo.FindByIDScoped(claims.SchoolCode, id)
	if err != nil || class.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Class not found"))
	}

	return dto.DetailClassResponse{
		ID: class.ID,
		ClassCode: dto.GeneralLabelKeyResponse{
			Key:   class.ClassCode,
			Label: class.DetailClassCode.Name,
		},
		ClassName: class.ClassName,
	}
}

// CreateNewClass creates a new class
func (s *ClassService) CreateNewClass(ctx *gin.Context, req dto.ModifyClassRequest) dto.DetailClassResponse {
	claims := jwt.GetDataClaims(ctx)

	newClass := school.Class{
		ClassCode:  req.ClassCode,
		ClassName:  req.ClassName,
		SchoolCode: claims.SchoolCode,
	}

	if err := s.repo.Create(&newClass); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	class, _ := s.repo.FindByIDScoped(claims.SchoolCode, newClass.ID)
	return dto.DetailClassResponse{
		ID: class.ID,
		ClassCode: dto.GeneralLabelKeyResponse{
			Key:   class.ClassCode,
			Label: class.DetailClassCode.Name,
		},
		ClassName: class.ClassName,
	}
}

// ModifyClass updates an existing class
func (s *ClassService) ModifyClass(ctx *gin.Context, id uint, req dto.ModifyClassRequest) dto.DetailClassResponse {
	claims := jwt.GetDataClaims(ctx)

	existingClass, err := s.repo.FindByIDScoped(claims.SchoolCode, id)
	if err != nil || existingClass.ID == 0 {
		panic(exception.NewNotFoundException(response.NotFound, "Class not found"))
	}

	existingClass.ClassName = req.ClassName
	existingClass.ClassCode = req.ClassCode

	if err := s.repo.Update(existingClass); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return dto.DetailClassResponse{
		ID: existingClass.ID,
		ClassCode: dto.GeneralLabelKeyResponse{
			Key:   existingClass.ClassCode,
			Label: existingClass.DetailClassCode.Name,
		},
		ClassName: existingClass.ClassName,
	}
}

// DeleteClass deletes a class
func (s *ClassService) DeleteClass(ctx *gin.Context, id uint) {
	claims := jwt.GetDataClaims(ctx)

	// Delete class members first
	s.repo.DeleteClassMembersByClassID(id)

	// Delete the class
	if err := s.repo.DeleteScoped(claims.SchoolCode, id); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
}

// MemberOfClass returns class members
func (s *ClassService) MemberOfClass(classID uint) []view.VStudent {
	members, _ := s.repo.FindMembersByClassID(classID)
	return members
}

// AddMemberOfClass adds students to a class
func (s *ClassService) AddMemberOfClass(req dto.ModifyClassMemberRequest) dto.ModifyClassMemberRequest {
	// Check for duplicates
	for _, studentID := range req.StudentId {
		member, _ := s.repo.FindStudentClassMembership(studentID, req.ClassId)
		if member != nil && member.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed add member, please check duplicate student added"))
		}
	}

	// Create members
	var members []student.StudentClass
	for _, studentID := range req.StudentId {
		members = append(members, student.StudentClass{
			StudentId: studentID,
			ClassId:   req.ClassId,
		})
	}

	if err := s.repo.CreateClassMembers(members); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return req
}

// RemoveMemberOfClass removes a class member
func (s *ClassService) RemoveMemberOfClass(id uint) {
	if err := s.repo.DeleteClassMemberByID(id); err != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Member not found"))
	}
}

// DeleteBatchMemberOfClass removes multiple class members
func (s *ClassService) DeleteBatchMemberOfClass(req dto.ModifyClassMemberRequest) dto.ModifyClassMemberRequest {
	for _, studentID := range req.StudentId {
		s.repo.DeleteClassMember(studentID, req.ClassId)
	}
	return req
}
