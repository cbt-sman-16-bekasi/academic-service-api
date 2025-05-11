package class_service

import (
	"errors"
	classRequest "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/class_request"
	response2 "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"
	classResponse "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/class_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
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
	paging := database.NewPagination[map[string]interface{}]().
		SetModal([]view.VClass{}).
		SetRequest(&request).
		FindAllPaging()
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

	if existingClass.ClassCode != request.ClassCode {
		checkCode := c.repoClass.GetClassByCode(request.ClassCode)
		if checkCode != nil {
			panic(exception.NewBadRequestExceptionStruct(response.ExpectationFailed, "Class code already exists"))
		}
	}

	existingClass.ClassName = request.ClassName
	c.repoClass.Database.Save(&existingClass)

	return classResponse.DetailClassResponse{}
}

func (c *ClassService) DeleteById(id uint) {
	c.repoClass.Delete(id)
}

func (c *ClassService) MemberOfClass(classId uint) []view.VStudent {
	var members []view.VStudent
	c.repoClass.Database.Where("class_id = ?", classId).Find(&members)
	return members
}

func (c *ClassService) AddMemberOfClass(request classRequest.ModifyClassMemberRequest) classRequest.ModifyClassMemberRequest {
	for _, member := range request.StudentId {
		var m student.StudentClass
		c.repoClass.Database.Where("student_id = ? AND class_id = ?", member, request.ClassId).First(&m)
		if m.ID != 0 {
			panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "Failed add member, please check duplicate student added"))
		}
	}

	var members []student.StudentClass
	for _, member := range request.StudentId {
		members = append(members, student.StudentClass{
			StudentId: member,
			ClassId:   request.ClassId,
		})
	}
	c.repoClass.Database.Create(&members)

	return request
}

func (c *ClassService) RemoveMemberOfClass(id uint) {
	if err := c.repoClass.Database.Where("id = ?", id).Delete(&student.StudentClass{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			panic(exception.NewBadRequestExceptionStruct(response.NotFound, "Failed remove member, please check duplicate student added"))
		}
		panic(exception.NewIntenalServerExceptionStruct(response.BadRequest, "Failed remove member, please try again!"))
	}
}
