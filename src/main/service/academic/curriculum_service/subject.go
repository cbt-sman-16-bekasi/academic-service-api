package curriculum_service

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/curriculum_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	response2 "github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
	"strings"
)

type SubjectService struct {
	repo *school_repository.SchoolRepository
}

func NewSubjectService() *SubjectService {
	return &SubjectService{repo: school_repository.NewSchoolRepository()}
}

func (s *SubjectService) CreateSubject(subject curriculum_request.SubjectRequest) curriculum_request.SubjectRequest {
	var existCode curriculum.Subject
	s.repo.Database.Where("code = ?", subject.Code).First(&existCode)
	if existCode.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Code already exists"))
	}

	data := curriculum.Subject{
		Model:   gorm.Model{},
		Code:    subject.Code,
		Subject: subject.Name,
	}
	s.repo.Database.Create(&data)
	return subject
}

func (s *SubjectService) UpdateSubject(id uint64, subject curriculum_request.SubjectRequest) curriculum_request.SubjectRequest {
	var existing curriculum.Subject
	s.repo.Database.Where("id = ?", id).First(&existing)
	if existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "ID not found"))
	}

	if strings.ToLower(existing.Code) != strings.ToLower(strings.Trim(subject.Code, " ")) {
		var existCode curriculum.Subject
		s.repo.Database.Where("code = ?", subject.Code).First(&existCode)
		if existCode.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "Code already exists"))
		}
	}

	existing.Subject = subject.Name
	existing.Code = subject.Code
	s.repo.Database.Save(&existing)

	return subject
}

func (s *SubjectService) DeleteSubject(id uint64) {
	var existCode curriculum.Subject
	s.repo.Database.Where("id = ?", id).First(&existCode)
	if existCode.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "ID not found"))
	}
	s.repo.Database.Delete(&existCode)
}

func (s *SubjectService) GetSubject(id uint64) curriculum.Subject {
	var existCode curriculum.Subject
	s.repo.Database.Where("id = ?", id).First(&existCode)
	if existCode.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response2.BadRequest, "ID not found"))
	}
	return existCode
}

func (s *SubjectService) GetAllSubject(request pagination.Request[map[string]interface{}]) *database.Paginator {
	return database.NewPagination[map[string]interface{}]().
		SetModal([]curriculum.Subject{}).
		SetRequest(&request).
		FindAllPaging()
}
