package class_service

import (
	classRequest "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/class_request"
	response2 "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"
	classResponse "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/class_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
)

type ClassService struct {
	repoClass *school_repository.ClassRepository
}

func NewClassService() *ClassService {
	return &ClassService{
		repoClass: school_repository.NewClassRepository(),
	}
}

func (c *ClassService) FindAllClass(request pagination.Request[map[string]interface{}]) *database.Paginator {
	paging := database.FindAllPaging(request, &school.Class{})
	return paging
}

func (c *ClassService) GetDetailClass(id uint) classResponse.DetailClassResponse {
	detail := c.repoClass.FindById(id)
	return classResponse.DetailClassResponse{
		ID: detail.ID,
		ClassCode: response2.GeneralLabelKeyResponse{
			Key:   detail.ClassCode,
			Label: detail.DetailClassCode.Name,
		},
		ClassName: detail.ClassName,
	}
}

func (c *ClassService) CreateNewClass(request classRequest.ModifyClassRequest) classResponse.DetailClassResponse {
	body := school.Class{
		ClassCode: request.ClassCode,
		ClassName: request.ClassName,
	}
	c.repoClass.Database.Create(&body)
	detail := c.repoClass.FindById(body.ID)
	return classResponse.DetailClassResponse{
		ID: detail.ID,
		ClassCode: response2.GeneralLabelKeyResponse{
			Key:   detail.ClassCode,
			Label: detail.DetailClassCode.Name,
		},
		ClassName: detail.ClassName,
	}
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

func (c *ClassService) DeleteById(id uint) {
	c.repoClass.Delete(id)
}
