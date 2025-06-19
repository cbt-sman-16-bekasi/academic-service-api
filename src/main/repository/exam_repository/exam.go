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
		Preload("ExamQuestion", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("ExamQuestion.QuestionOption").Preload(clause.Associations).First(&exam)
	return &exam
}

func (e *ExamRepository) FindByIdQuestion(id uint) *school.ExamQuestion {
	var question school.ExamQuestion
	e.Database.Where("id = ?", id).Preload("QuestionOption", func(db *gorm.DB) *gorm.DB { return db.Order("answer_id asc") }).First(&question)

	return &question
}

func (e *ExamRepository) GetExamData(classId uint) []school.Exam {
	var examMember []school.ExamMember
	e.Database.Where("class = ?", classId).Preload("DetailExam").Find(&examMember)

	var examCode []string
	for _, examM := range examMember {
		examCode = append(examCode, examM.ExamCode)
	}

	var exam []school.Exam
	e.Database.Debug().Where("code IN ?", examCode).Find(&exam)

	return exam
}

func (e *ExamRepository) GetExamSessionActiveNow(examCode []string, studentId uint) *school.ExamSession {
	var examSession school.ExamSession
	timeNow := time.Now()
	e.Database.Where("exam_session.session_id IN ?", examCode).
		Joins("LEFT JOIN cbt_service.student_history_taken st ON st.session_id = exam_session.session_id AND st.student_id = ? AND status != 'COMPLETED'", studentId).
		Where("start_date <= ? and end_date >= ?", timeNow, timeNow).
		Order("exam_session.start_date asc").
		First(&examSession)

	return &examSession
}

func (e *ExamRepository) FindByCode(code string) school.Exam {
	var exam school.Exam
	e.Database.Where("code = ?", code).
		Preload("DetailSubject").
		Preload("ExamMember").
		Preload("ExamMember.DetailClass").
		Preload("ExamMember.DetailClass.DetailClassCode").
		Preload("DetailTypeExam").
		Preload("DetailTypeExam.DetailRole").
		Preload("ExamQuestion").
		Preload("ExamQuestion.QuestionOption").Preload(clause.Associations).
		First(&exam)
	return exam
}
