package exam_service

import (
	"fmt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/exam_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/exam_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
)

type TypeExamService struct {
	typeExamRepo *exam_repository.TypeExamRepository
}

func NewTypeExamService() *TypeExamService {
	return &TypeExamService{
		typeExamRepo: exam_repository.NewTypeExamRepository(),
	}
}

func (t *TypeExamService) GetAll(request pagination.Request[map[string]interface{}]) *database.Paginator {
	return database.NewPagination[map[string]interface{}]().
		SetModal([]school.TypeExam{}).
		SetRequest(&request).
		SetPreloads("DetailRole").FindAllPaging()
}

func (t *TypeExamService) GetDetail(id uint) *school.TypeExam {
	typeExam := t.typeExamRepo.Repository.FindById(id)
	if typeExam.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Type exam with id '%d' not found", id)))
	}
	return typeExam
}

func (t *TypeExamService) CreateTypeExam(request exam_request.ModifyTypeExamRequest) school.TypeExam {
	checkCode := t.typeExamRepo.FindByCode(request.CodeTypeExam)
	if checkCode != nil {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, fmt.Sprintf("Type exam code '%s' already exist", request.CodeTypeExam)))
	}

	newData := school.TypeExam{
		Code:  request.CodeTypeExam,
		Name:  request.TypeExam,
		Color: request.Color,
		Role:  request.Role,
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
	_ = t.typeExamRepo.Repository.DeleteById(id)
}
