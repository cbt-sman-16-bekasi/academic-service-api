package student_service

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/request/student_request"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/student_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/pagination"
)

type StudentService struct {
	studentRepo *school_repository.StudentRepository
}

func NewStudentService() *StudentService {
	return &StudentService{
		studentRepo: school_repository.NewStudentRepository(),
	}
}

func (s *StudentService) DetailStudent(id uint) student_response.DetailStudentResponse {
	detail := s.studentRepo.FindById(id)
	return student_response.DetailStudentResponse{
		Nisn:   detail.DetailStudent.DetailUser.Username,
		Name:   detail.DetailStudent.Name,
		Gender: detail.DetailStudent.Gender,
		Class: response.GeneralLabelKeyResponse{
			Key:   detail.DetailClass.ClassCode,
			Label: detail.DetailClass.ClassName,
		},
	}
}

func (s *StudentService) AllStudent(request pagination.Request[map[string]interface{}]) *database.Paginator {
	paging := database.NewPagination[map[string]interface{}]().
		SetRequest(&request).
		SetModal(&student.StudentClass{}).
		SetPreloads("DetailStudent", "DetailClass", "DetailStudent.DetailUser").FindAllPaging()
	return paging
}

func (s *StudentService) CreateStudent(request student_request.StudentModifyRequest) student_response.DetailStudentResponse {
	return student_response.DetailStudentResponse{}
}

func (s *StudentService) UpdateStudent(id uint, request student_request.StudentModifyRequest) student_response.DetailStudentResponse {
	return student_response.DetailStudentResponse{}
}

func (s *StudentService) DeleteById(id uint) {
	s.studentRepo.Delete(id)
}
