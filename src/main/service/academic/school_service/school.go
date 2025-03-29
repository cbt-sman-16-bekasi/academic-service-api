package school_service

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/class_request"
	response2 "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"
	classResponse "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/class_response"
	schoolResponse "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/school_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/gin-gonic/gin"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/pagination"
	"github.com/yon-module/yon-framework/server/response"
)

type SchoolService struct {
	repo             *school_repository.SchoolRepository
	repoClassSubject *school_repository.ClassSubjectRepository
	repoClassCode    *school_repository.ClassCodeRepository
}

func NewSchoolService() *SchoolService {
	return &SchoolService{
		repo:             school_repository.NewSchoolRepository(),
		repoClassSubject: school_repository.NewClassSubjectRepository(),
		repoClassCode:    school_repository.NewClassCodeRepository(),
	}
}

func (s *SchoolService) GetAllClassCode() []classResponse.ClassCodeResponse {
	classCode := s.repoClassCode.GetAllClassCode()
	var classes []classResponse.ClassCodeResponse
	for _, class := range classCode {
		classes = append(classes, classResponse.ClassCodeResponse{
			Code: class.Code,
			Name: class.Name,
		})
	}
	return classes
}

func (s *SchoolService) GetAllSubject() []curriculum.Subject {
	subjects := s.repo.AllSubject()
	return subjects
}

func (s *SchoolService) RetrieveDetailSchool(c *gin.Context) schoolResponse.DetailSchool {
	schoolDetail := s.repo.FindTopBySchoolCode(c.Query("schoolCode"))
	return schoolResponse.DetailSchool{
		SchoolName: schoolDetail.SchoolName,
		Id:         schoolDetail.SchoolCode,
		Address:    schoolDetail.Address,
		Email:      schoolDetail.Email,
		Phone:      schoolDetail.Phone,
		Logo:       schoolDetail.Logo,
	}
}

func (s *SchoolService) GetAllClassSubject(request pagination.Request[map[string]interface{}]) *database.Paginator {
	page := database.FindAllPaging(request, []school.ClassSubject{})
	return page
}

func (s *SchoolService) GetDetailClassSubject(id uint) classResponse.DetailClassSubjectResponse {
	detail := s.repoClassSubject.FindById(id)

	return classResponse.DetailClassSubjectResponse{
		ID: detail.ID,
		ClassCode: response2.GeneralLabelKeyResponse{
			Key:   detail.ClassCode,
			Label: detail.DetailClassCode.Name,
		},
		Subject: response2.GeneralLabelKeyResponse{
			Key:   detail.SubjectCode,
			Label: detail.DetailSubject.Subject,
		},
	}
}

func (s *SchoolService) CreateClassSubject(request class_request.ModifyClassSubject) classResponse.DetailClassSubjectResponse {
	var existing *school.ClassSubject
	s.repoClassSubject.Database.Where("class_code = ? AND subject_code = ?", request.ClassCode, request.SubjectCode).First(&existing)

	if existing.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "This class has already been created"))
	}

	newData := &school.ClassSubject{
		SubjectCode: request.SubjectCode,
		ClassCode:   request.ClassCode,
	}

	s.repoClassSubject.Database.Create(newData)

	return classResponse.DetailClassSubjectResponse{
		ID:        newData.ID,
		ClassCode: response2.GeneralLabelKeyResponse{},
		Subject:   response2.GeneralLabelKeyResponse{},
	}
}

func (s *SchoolService) UpdateClassSubject(id uint, request class_request.ModifyClassSubject) classResponse.DetailClassSubjectResponse {
	existing := s.repoClassSubject.FindById(id)
	if existing.ID == 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "This class not found"))
	}

	var duplicateCheck *school.ClassSubject
	s.repoClassSubject.Database.Where("class_code = ? AND subject_code = ?", request.ClassCode, request.SubjectCode).First(&duplicateCheck)

	if duplicateCheck.ID != 0 {
		panic(exception.NewBadRequestExceptionStruct(response.BadRequest, "This class has already been created"))
	}

	if existing.ClassCode != request.ClassCode {
		existing.ClassCode = request.ClassCode
	}

	if existing.SubjectCode != request.SubjectCode {
		existing.SubjectCode = request.SubjectCode
	}

	s.repoClassSubject.Database.Save(&existing)

	return classResponse.DetailClassSubjectResponse{
		ID:        id,
		ClassCode: response2.GeneralLabelKeyResponse{},
		Subject:   response2.GeneralLabelKeyResponse{},
	}
}

func (s *SchoolService) DeleteClassSubject(id uint) {
	s.repoClassSubject.DeleteById(id)
}

func (s *SchoolService) DashboardUser() schoolResponse.DashboardResponse {
	return schoolResponse.DashboardResponse{}
}
