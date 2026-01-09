package service

import (
	"fmt"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/helper/jwt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/exam/repository"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/exception"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/pagination"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/gin-gonic/gin"
)

type TypeExamService struct {
	typeExamRepo *repository.TypeExamRepository
}

// NewTypeExamService creates service with injected repository
func NewTypeExamService(repo *repository.TypeExamRepository) *TypeExamService {
	return &TypeExamService{
		typeExamRepo: repo,
	}
}

func (t *TypeExamService) GetAll(c *gin.Context, request pagination.Request[map[string]interface{}]) *database.Paginator {
	claims := jwt.GetDataClaims(c)
	filter := map[string]interface{}{}

	if request.Filter != nil {
		filter = *request.Filter
	}

	if claims.Role != "ADMIN" {
		filter["role"] = claims.Role
	}

	filter["school_code"] = claims.SchoolCode
	request.Filter = &filter

	return database.NewPagination[map[string]interface{}](t.typeExamRepo.Database).
		SetModel([]school.TypeExam{}).
		SetRequest(&request).
		SetPreloads("DetailRole").FindAllPaging()
}

func (t *TypeExamService) GetDetail(id uint) *school.TypeExam {
	typeExam := t.typeExamRepo.FindById(id)
	if typeExam == nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Type exam with id '%d' not found", id)))
	}
	return typeExam
}

func (t *TypeExamService) CreateTypeExam(c *gin.Context, request exam_request.ModifyTypeExamRequest) school.TypeExam {
	claims := jwt.GetDataClaims(c)
	checkCode := t.typeExamRepo.FindByCode(request.CodeTypeExam)
	if checkCode != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Type exam code '%s' already exist", request.CodeTypeExam)))
	}

	newData := school.TypeExam{
		Code:       request.CodeTypeExam,
		Name:       request.TypeExam,
		Color:      request.Color,
		Role:       request.Role,
		SchoolCode: claims.SchoolCode,
	}

	t.typeExamRepo.Database.Create(&newData)
	return newData
}

func (t *TypeExamService) ModifyTypeExam(id uint, request exam_request.ModifyTypeExamRequest) *school.TypeExam {
	existingData := t.typeExamRepo.Repository.FindById(id)
	if existingData.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Type exam id '%d' not found", id)))
	}

	if existingData.Code != request.CodeTypeExam {
		checkCode := t.typeExamRepo.FindByCode(request.CodeTypeExam)
		if checkCode != nil {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Type exam code '%s' already exist", request.CodeTypeExam)))
		}

		existingData.Code = request.CodeTypeExam
	}

	existingData.Name = request.TypeExam
	existingData.Color = request.Color
	existingData.Role = request.Role
	t.typeExamRepo.Database.Save(&existingData)

	return existingData
}

func (t *TypeExamService) DeleteTypeExam(id uint) {
	existingData := t.typeExamRepo.Repository.FindById(id)
	if existingData.ID == 0 {
		panic("Failed delete type exam")
	}
	t.typeExamRepo.Database.Delete(&existingData)
}
