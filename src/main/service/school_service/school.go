package school_service

import (
	schoolResponse "github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/dto/response/school_response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/school"
	"github.com/gin-gonic/gin"
)

type SchoolService struct {
	repo *school.SchoolRepository
}

func NewSchoolService() *SchoolService {
	return &SchoolService{repo: school.NewSchoolRepository()}
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
