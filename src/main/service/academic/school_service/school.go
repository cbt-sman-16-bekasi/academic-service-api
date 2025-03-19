package school_service

import (
	schoolResponse "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/school_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school_repository"
	"github.com/gin-gonic/gin"
)

type SchoolService struct {
	repo *school_repository.SchoolRepository
}

func NewSchoolService() *SchoolService {
	return &SchoolService{repo: school_repository.NewSchoolRepository()}
}

func (s SchoolService) RetrieveDetailSchool(c *gin.Context) schoolResponse.DetailSchool {
	school := s.repo.FindTopBySchoolCode(c.Query("schoolCode"))
	return schoolResponse.DetailSchool{
		SchoolName: school.SchoolName,
		Id:         school.SchoolCode,
		Address:    school.Address,
		Email:      school.Email,
		Phone:      school.Phone,
		Logo:       school.Logo,
	}
}
