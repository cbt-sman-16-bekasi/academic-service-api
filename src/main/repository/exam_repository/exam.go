package exam_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ExamRepository struct {
	Database   *gorm.DB
	Repository *database.GpaRepository[school.Exam]
}

func NewExamRepository() *ExamRepository {
	return &ExamRepository{
		Database:   database.GetDB(),
		Repository: database.NewGpaRepository(school.Exam{}),
	}
}

func (e *ExamRepository) FindById(id uint) *school.Exam {
	var exam school.Exam
	e.Database.Debug().Where("id = ?", id).
		Preload("DetailSubject").
		Preload("ExamMember").
		Preload("ExamMember.DetailClass").
		Preload("ExamMember.DetailClass.DetailClassCode").
		Preload("DetailTypeExam").
		Preload("DetailTypeExam.DetailRole").
		Preload("ExamQuestion").
		Preload("ExamQuestion.QuestionOption").Preload(clause.Associations).First(&exam)
	return &exam
}

func (e *ExamRepository) FindByIdQuestion(id uint) *school.ExamQuestion {
	var question school.ExamQuestion
	e.Database.Where("id = ?", id).Preload("QuestionOption").First(&question)

	return &question
}

func (e *ExamRepository) GetExamData(classId uint) *school.Exam {
	var examMember school.ExamMember
	e.Database.Where("class = ?", classId).Preload("DetailExam").First(&examMember)

	return &examMember.DetailExam
}

func (e *ExamRepository) GetExamSessionActiveNow(examCode string) *school.ExamSession {
	var examSession school.ExamSession
	timeNow := time.Now()
	e.Database.Where("exam_code = ?", examCode).
		Where("start_date <= ? and end_date >= ?", timeNow, timeNow).
		First(&examSession)

	return &examSession
}

func (e *ExamRepository) FindByCode(code string) school.Exam {
	var exam school.Exam
	e.Database.Where("code = ?", code).
		Preload("ExamMember").
		Preload("ExamMember.DetailClass").
		First(&exam)
	return exam
}
