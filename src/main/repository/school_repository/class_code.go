package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"github.com/yon-module/yon-framework/exception"
	"github.com/yon-module/yon-framework/server/response"
	"gorm.io/gorm"
)

type ClassCodeRepository struct {
	Database *gorm.DB
}

func NewClassCodeRepository() *ClassCodeRepository {
	return &ClassCodeRepository{
		Database: database.GetDB(),
	}
}

func (repo *ClassCodeRepository) GetAllClassCode() []school.ClassCode {
	var classes []school.ClassCode
	q := repo.Database.Find(&classes)
	if q.Error != nil {
		panic(exception.NewIntenalServerExceptionStruct(response.ServerError, "Data class code not found"))
	}
	return classes
}
