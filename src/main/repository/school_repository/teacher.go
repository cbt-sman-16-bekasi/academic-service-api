package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type TeacherRepository struct {
	Database *gorm.DB
}

func NewTeacherRepository() *TeacherRepository {
	return &TeacherRepository{
		Database: database.GetDB(),
	}
}

func (s *TeacherRepository) Delete(id uint) {
	s.Database.Where("id = ?", id).Delete(&teacher.Teacher{})
}

func (s *TeacherRepository) FindById(id uint) *teacher.Teacher {
	var teach teacher.Teacher
	s.Database.Where("id = ?", id).Preload("DetailUser").Preload("DetailUser.RoleUser").First(&teach)

	if teach.ID == 0 {
		return nil
	}
	return &teach
}
