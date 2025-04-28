package exam_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type ExamSessionRepository struct {
	Database   *gorm.DB
	Repository *database.GpaRepository[school.ExamSession]
}

func NewExamSessionRepository() *ExamSessionRepository {
	return &ExamSessionRepository{
		Database:   database.GetDB(),
		Repository: database.NewGpaRepository(school.ExamSession{}),
	}
}

func (e *ExamSessionRepository) FindById(id uint) *school.ExamSession {
	var session school.ExamSession
	e.Database.Where("id = ?", id).Preload("DetailExam").
		Preload("DetailExam.DetailSubject").
		Preload("DetailExam.DetailTypeExam").
		Preload("ExamSessionMember").
		Preload("ExamSessionMember.DetailClass").
		First(&session)
	return &session
}
