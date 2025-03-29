package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
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

func (c *ClassSubjectRepository) FindById(id uint) (classSubject *school.ClassSubject) {
	c.Database.Where("id = ?", id).Preload("DetailSubject").Preload("DetailClassCode").First(&classSubject)
	return
}

func (c *ClassSubjectRepository) DeleteById(id uint) {
	c.Database.Where("id = ?", id).Delete(&school.ClassSubject{})
}
