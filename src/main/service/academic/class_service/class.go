package class_service

import (
	classRequest "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/class_request"
	classResponse "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/class_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
)

type ClassService struct {
	repoClass        *school_repository.ClassRepository
	repoClassCode    *school_repository.ClassCodeRepository
	repoClassSubject *school_repository.ClassSubjectRepository
}

func NewClassService() *ClassService {
	return &ClassService{
		repoClass:        school_repository.NewClassRepository(),
		repoClassCode:    school_repository.NewClassCodeRepository(),
		repoClassSubject: school_repository.NewClassSubjectRepository(),
	}
}

func (c *ClassService) GetAllClassCode() []classResponse.ClassCodeResponse {
	classCode := c.repoClassCode.GetAllClassCode()
	var classes []classResponse.ClassCodeResponse
	for _, class := range classCode {
		classes = append(classes, classResponse.ClassCodeResponse{
			Code: class.Code,
			Name: class.Name,
		})
	}
	return classes
}

func (c *ClassService) CreateNewClass(request classRequest.ModifyClassRequest) classResponse.DetailClassResponse {
	existingClass := c.repoClass.GetClassByCode(request.ClassCode)
	if existingClass != nil && existingClass.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Class with this code already exists"))
	}

	body := school.Class{
		Model:     gorm.Model{},
		ClassCode: request.ClassCode,
		ClassName: request.ClassName,
	}
	c.repoClass.Database.Create(&body)

	return classResponse.DetailClassResponse{}
}

func (c *ClassService) ModifyClass(id uint, request classRequest.ModifyClassRequest) classResponse.DetailClassResponse {
	existingClass := c.repoClass.FindById(id)
	if existingClass == nil || existingClass.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.NotFound, "Class with this code not found"))
	}

	checkCode := c.repoClass.GetClassByCode(request.ClassCode)
	if checkCode.ID != 0 && checkCode.ID != id {
		panic(exception.NewBadRequestExceptionStruct(response.ExpectationFailed, "Class code already exists"))
	}

	existingClass.ClassName = request.ClassName
	c.repoClass.Database.Save(&existingClass)

	return classResponse.DetailClassResponse{}
}

func (c *ClassService) GetAllClassSubject(request pagination.Request[map[string]interface{}]) *database.Paginator {
	page := database.FindAllPaging(request, &school.ClassSubject{})
	return page
}
