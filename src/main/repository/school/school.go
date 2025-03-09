package school

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/server/response"
)

type SchoolRepository struct {
}

func NewSchoolRepository() *SchoolRepository {
	return &SchoolRepository{}
}

func (s SchoolRepository) FindTopBySchoolCode(schoolCode string) *school.School {
	var school school.School
	q := database.GetDB().Model(&school).Where("school_code=?", schoolCode).First(&school)
	if q.Error != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Data school not found"))
	}
	return &school
}
