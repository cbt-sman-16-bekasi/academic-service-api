package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type ClassRepository struct {
	Database *gorm.DB
}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{
		Database: database.GetDB(),
	}
}

func (repo *ClassRepository) GetClassByCode(classCode string) *school.Class {
	var class school.Class
	_ = repo.Database.Where("class_code = ?", classCode).First(&class)
	return &class
}

func (repo *ClassRepository) FindById(id uint) *school.Class {
	var class school.Class
	_ = repo.Database.Where("id = ?", id).First(&class)
	return &class
}

func (repo *ClassRepository) AllClass() []school.Class {
	var classes []school.Class
	_ = repo.Database.Preload("DetailClassCode").Find(&classes)
	return classes
}
