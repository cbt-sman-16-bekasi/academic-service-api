package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
)

type SchoolRepository struct {
	Database *gorm.DB
}

func NewSchoolRepository() *SchoolRepository {
	return &SchoolRepository{
		Database: database.GetDB(),
	}
}

func (s SchoolRepository) FindTopBySchoolCode(schoolCode string) *school.School {
	var school school.School
	q := s.Database.Where("school_code=?", schoolCode).First(&school)
	if q.Error != nil {
		panic(exception.NewNotFoundException(response.NotFound, "Data school_repository not found"))
	}
	return &school
}
