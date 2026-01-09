package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"gorm.io/gorm"
)

type TypeExamRepository struct {
	Database   *gorm.DB
	Repository *database.GpaRepository[school.TypeExam]
}

// NewTypeExamRepository creates repository with injected DB
func NewTypeExamRepository(db *gorm.DB) *TypeExamRepository {
	return &TypeExamRepository{
		Database:   db,
		Repository: database.NewGpaRepository(school.TypeExam{}, db),
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

// FindByCodeScoped returns type exam by code filtered by school_code
func (t *TypeExamRepository) FindByCodeScoped(schoolCode, code string) *school.TypeExam {
	var typeExam school.TypeExam
	t.Database.Scopes(database.SchoolScope(schoolCode)).
		Where("code = ?", code).Preload("DetailRole").Find(&typeExam)

	if typeExam.ID == 0 {
		return nil
	}
	return &typeExam
}

func (t *TypeExamRepository) FindById(id uint) *school.TypeExam {
	var typeExam school.TypeExam
	t.Database.Where("id = ?", id).Preload("DetailRole").Find(&typeExam)

	if typeExam.ID == 0 {
		return nil
	}
	return &typeExam
}

// FindByIdScoped returns type exam by id filtered by school_code
func (t *TypeExamRepository) FindByIdScoped(schoolCode string, id uint) *school.TypeExam {
	var typeExam school.TypeExam
	t.Database.Scopes(database.SchoolScope(schoolCode)).
		Where("id = ?", id).Preload("DetailRole").Find(&typeExam)

	if typeExam.ID == 0 {
		return nil
	}
	return &typeExam
}

// AllTypeExamScoped returns all type exams filtered by school_code
func (t *TypeExamRepository) AllTypeExamScoped(schoolCode string) []school.TypeExam {
	var typeExams []school.TypeExam
	t.Database.Scopes(database.SchoolScope(schoolCode)).
		Preload("DetailRole").Find(&typeExams)
	return typeExams
}
