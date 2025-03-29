package exam_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type TypeExamRepository struct {
	Database *gorm.DB
}

func NewTypeExamRepository() *TypeExamRepository {
	return &TypeExamRepository{
		Database: database.GetDB(),
	}
}

func (t *TypeExamRepository) FindByCode(code string) *school.TypeExam {
	var typeExam school.TypeExam
	t.Database.Where("code = ?", code).Preload("DetailRole").Find(&typeExam)
	return &typeExam
}
