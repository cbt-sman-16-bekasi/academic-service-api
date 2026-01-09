package service

import (
	"strings"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/dto"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/entity"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/gin-gonic/gin"

	// Using existing school entity for ClassSubject (cross-module dependency)
	schoolEntity "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
)

type SubjectService struct {
	repo *repository.SubjectRepository
}

// NewSubjectService creates service with injected repository
func NewSubjectService(repo *repository.SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

// GetAllSubject returns paginated subjects
func (s *SubjectService) GetAllSubject(ctx *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(ctx)

	filter := map[string]interface{}{}
	if request.Filter != nil {
		filter = *request.Filter
	}
	filter["school_code"] = claims.SchoolCode
	request.Filter = &filter

	return database.NewPagination[map[string]interface{}](s.repo.DB()).
		SetModel([]entity.VSubject{}).
		SetRequest(&request).
		FindAllPaging()
}

// GetSubject returns subject by ID
func (s *SubjectService) GetSubject(id uint64) *entity.VSubject {
	subject, err := s.repo.FindViewByID(uint(id))
	if err != nil || subject.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.NotFound, "Subject not found"))
	}
	return subject
}

// CreateSubject creates a new subject
func (s *SubjectService) CreateSubject(ctx *gin.Context, req dto.SubjectRequest) dto.SubjectRequest {
	// Check if code exists
	claims := jwt.GetDataClaims(ctx)
	schoolCode := claims.SchoolCode
	if s.repo.Exists(schoolCode, req.Code) {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Code already exists"))
	}

	tx := s.repo.DB().Begin()

	subject := entity.Subject{
		Code:       req.Code,
		Subject:    req.Name,
		SchoolCode: schoolCode,
	}

	if err := tx.Create(&subject).Error; err != nil {
		tx.Rollback()
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	// Delete existing class-subject relations
	if err := tx.Where("subject_code = ? and school_code = ?", subject.Code, schoolCode).Delete(&schoolEntity.ClassSubject{}).Error; err != nil {
		tx.Rollback()
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	// Create new class-subject relations
	if len(req.ClassCode) > 0 {
		var classSubjects []schoolEntity.ClassSubject
		for _, c := range req.ClassCode {
			classSubjects = append(classSubjects, schoolEntity.ClassSubject{
				ClassCode:   c,
				SubjectCode: subject.Code,
				SchoolCode:  schoolCode,
			})
		}

		if err := tx.Create(&classSubjects).Error; err != nil {
			tx.Rollback()
			panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return req
}

// UpdateSubject updates an existing subject
func (s *SubjectService) UpdateSubject(ctx *gin.Context, id uint64, req dto.SubjectRequest) dto.SubjectRequest {
	claims := jwt.GetDataClaims(ctx)
	schoolCode := claims.SchoolCode

	existing, err := s.repo.FindByID(uint(id))
	if err != nil || existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.NotFound, "Subject not found"))
	}

	// Check if new code conflicts with another subject
	if strings.ToLower(existing.Code) != strings.ToLower(strings.TrimSpace(req.Code)) {
		if s.repo.ExistsExcluding(req.Code, existing.ID) {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Code already exists"))
		}
	}

	tx := s.repo.DB().Begin()

	existing.Subject = req.Name
	existing.Code = req.Code

	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	// Update class-subject relations if provided
	if req.ClassCode != nil && len(req.ClassCode) > 0 {
		if err := tx.Where("subject_code = ? and school_code = ?", existing.Code, schoolCode).Delete(&schoolEntity.ClassSubject{}).Error; err != nil {
			tx.Rollback()
			panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
		}

		var classSubjects []schoolEntity.ClassSubject
		for _, c := range req.ClassCode {
			classSubjects = append(classSubjects, schoolEntity.ClassSubject{
				ClassCode:   c,
				SubjectCode: existing.Code,
				SchoolCode:  schoolCode,
			})
		}

		if err := tx.Create(&classSubjects).Error; err != nil {
			tx.Rollback()
			panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}

	return req
}

// DeleteSubject deletes a subject
func (s *SubjectService) DeleteSubject(id uint64) {
	existing, err := s.repo.FindByID(uint(id))
	if err != nil || existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.NotFound, "Subject not found"))
	}

	if err := s.repo.Delete(existing.ID); err != nil {
		panic(exception.NewInternalServerExceptionStruct(response.ServerError, err.Error()))
	}
}
