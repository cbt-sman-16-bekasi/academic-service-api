package exam_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type ExamRepository struct {
	Database *gorm.DB
}

func NewExamRepository() *ExamRepository {
	return &ExamRepository{
		Database: database.GetDB(),
	}
}

func (e *ExamRepository) FindById(id uint) *school.Exam {
	var exam school.Exam
	e.Database.Where("id = ?", id).Preload("DetailSubject").Preload("ExamMember").First(&exam)
	return &exam
}
