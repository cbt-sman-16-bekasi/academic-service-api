package exam_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type TypeExamRepository struct {
	Database   *gorm.DB
	Repository *database.GpaRepository[school.TypeExam]
}

func NewTypeExamRepository() *TypeExamRepository {
	database.NewGpaRepository(school.TypeExam{})
	return &TypeExamRepository{
		Database:   database.GetDB(),
		Repository: database.NewGpaRepository(school.TypeExam{}),
	}
}

func (t *TypeExamRepository) FindByCode(code string) *school.TypeExam {
	var typeExam school.TypeExam
	t.Database.Where("code = ?", code).Preload("DetailRole").Find(&typeExam)

	if typeExam.ID == 0 {
		return nil
	}
	return &typeExam
}
