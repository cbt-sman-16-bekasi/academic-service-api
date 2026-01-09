package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"gorm.io/gorm"
)

type ExamSessionRepository struct {
	Database   *gorm.DB
	Repository *database.GpaRepository[school.ExamSession]
}

// NewExamSessionRepository creates repository with injected DB
func NewExamSessionRepository(db *gorm.DB) *ExamSessionRepository {
	return &ExamSessionRepository{
		Database:   db,
		Repository: database.NewGpaRepository(school.ExamSession{}, db),
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

// FindByIdScoped returns exam session by id filtered by school_code
func (e *ExamSessionRepository) FindByIdScoped(schoolCode string, id uint) *school.ExamSession {
	var session school.ExamSession
	e.Database.Scopes(database.SchoolScope(schoolCode)).
		Where("id = ?", id).Preload("DetailExam").
		Preload("DetailExam.DetailSubject").
		Preload("DetailExam.DetailTypeExam").
		Preload("ExamSessionMember").
		Preload("ExamSessionMember.DetailClass").
		First(&session)
	return &session
}

// FindBySessionIdScoped returns exam session by session_id filtered by school_code
func (e *ExamSessionRepository) FindBySessionIdScoped(schoolCode, sessionId string) *school.ExamSession {
	var session school.ExamSession
	e.Database.Scopes(database.SchoolScope(schoolCode)).
		Where("session_id = ?", sessionId).Preload("DetailExam").
		Preload("DetailExam.DetailSubject").
		Preload("DetailExam.DetailTypeExam").
		Preload("ExamSessionMember").
		Preload("ExamSessionMember.DetailClass").
		First(&session)
	return &session
}
