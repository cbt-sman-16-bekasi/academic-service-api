package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/repository/exam_repository"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StudentRepository struct {
	ExamRepo *exam_repository.ExamRepository
	Database *gorm.DB
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		Database: database.GetDB(),
		ExamRepo: exam_repository.NewExamRepository(),
	}
}

func (s *StudentRepository) FindById(id uint) (student *student.StudentClass) {
	_ = s.Database.Where("id = ?", id).
		Preload("DetailClass").
		Preload("DetailStudent.DetailUser").
		Preload("DetailStudent.DetailUser.RoleUser").
		Preload("DetailStudent").
		Preload(clause.Associations).
		First(&student)
	return
}

func (s *StudentRepository) FindByNISN(nisn string) (student *student.Student) {
	_ = s.Database.Where("nisn = ?", nisn).
		Preload("DetailUser").Preload("DetailUser.RoleUser").
		Preload(clause.Associations).
		First(&student)
	return
}

func (s *StudentRepository) GetStudentClass(studentId uint) (student *student.StudentClass) {
	_ = s.Database.Where("student_id = ?", studentId).
		Preload("DetailClass").
		Preload("DetailStudent.DetailUser").
		Preload("DetailStudent.DetailUser.RoleUser").
		Preload("DetailStudent").
		Preload(clause.Associations).
		First(&student)
	return
}

func (s *StudentRepository) Delete(id uint) {
	s.Database.Where("id = ?", id).Delete(&student.StudentClass{})
}
