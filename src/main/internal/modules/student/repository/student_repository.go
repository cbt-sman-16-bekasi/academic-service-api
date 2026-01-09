package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

// NewStudentRepository creates repository with injected DB
func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// DB returns the database connection for complex queries
func (r *StudentRepository) DB() *gorm.DB {
	return r.db
}

// FindByID finds student by ID (from view)
func (r *StudentRepository) FindByID(id uint) (*view.VStudent, error) {
	var s view.VStudent
	err := r.db.Where("id = ?", id).First(&s).Error
	return &s, err
}

// FindByNISN finds student by NISN
func (r *StudentRepository) FindByNISN(nisn string) (*student.Student, error) {
	var s student.Student
	err := r.db.Where("nisn = ?", nisn).First(&s).Error
	return &s, err
}

// FindStudentEntityByID finds student entity by ID
func (r *StudentRepository) FindStudentEntityByID(id uint) (*student.Student, error) {
	var s student.Student
	err := r.db.Where("id = ?", id).First(&s).Error
	return &s, err
}

// FindStudentClassByStudentID finds student class by student ID
func (r *StudentRepository) FindStudentClassByStudentID(studentID uint) (*student.StudentClass, error) {
	var sc student.StudentClass
	err := r.db.Where("student_id = ?", studentID).
		Preload("DetailStudent.DetailUser").
		Preload("DetailStudent").
		First(&sc).Error
	return &sc, err
}

// Create creates a new student
func (r *StudentRepository) Create(s *student.Student) error {
	return r.db.Create(s).Error
}

// Update updates a student
func (r *StudentRepository) Update(id uint, s *student.Student) error {
	return r.db.Model(&student.Student{}).Where("id = ?", id).Updates(s).Error
}

// Delete deletes a student and related records
func (r *StudentRepository) Delete(id uint) error {
	return r.db.Where("student_id = ?", id).Delete(&student.StudentClass{}).Error
}

// DeleteStudent deletes student entity
func (r *StudentRepository) DeleteStudent(id uint) error {
	return r.db.Where("id = ?", id).Delete(&student.Student{}).Error
}

// DeleteUser deletes user entity
func (r *StudentRepository) DeleteUser(id uint) error {
	return r.db.Where("id = ?", id).Delete(&user.User{}).Error
}

// CreateStudentClass creates a student class record
func (r *StudentRepository) CreateStudentClass(sc *student.StudentClass) error {
	return r.db.Create(sc).Error
}

// UpdateStudentClass updates student class
func (r *StudentRepository) UpdateStudentClass(studentID uint, classID uint) error {
	return r.db.Model(&student.StudentClass{}).Where("student_id = ?", studentID).Update("class_id", classID).Error
}
