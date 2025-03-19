package school_repository

import (
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type ClassSubjectRepository struct {
	Database *gorm.DB
}

func NewClassSubjectRepository() *ClassSubjectRepository {
	return &ClassSubjectRepository{
		Database: database.GetDB(),
	}
}
