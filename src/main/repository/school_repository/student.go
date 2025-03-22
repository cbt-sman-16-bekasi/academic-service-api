package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type StudentRepository struct {
	Database *gorm.DB
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		Database: database.GetDB(),
	}
}

func (s *StudentRepository) FindById(id uint) (student *student.StudentClass) {
	_ = s.Database.Where("id = ?", id).
		Preload("DetailStudent", "DetailClass", "DetailStudent.DetailUser").
		First(&student)
	return
}

func (s *StudentRepository) Delete(id uint) {
	s.Database.Where("id = ?", id).Delete(&school.Class{})
}
