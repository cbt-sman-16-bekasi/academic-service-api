package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"gorm.io/gorm"
)

type TeacherRepository struct {
	db *gorm.DB
}

// NewTeacherRepository creates repository with injected DB
func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

// DB returns the database connection for complex queries
func (r *TeacherRepository) DB() *gorm.DB {
	return r.db
}

// FindByID finds teacher by ID with preloads
func (r *TeacherRepository) FindByID(id uint) (*teacher.Teacher, error) {
	var t teacher.Teacher
	err := r.db.Where("id = ?", id).
		Preload("DetailUser").
		Preload("DetailUser.RoleUser").
		First(&t).Error
	return &t, err
}

// FindByNuptk finds teacher by NUPTK
func (r *TeacherRepository) FindByNuptk(schoolCode, nuptk string) (*teacher.Teacher, error) {
	var t teacher.Teacher
	err := r.db.Where("nuptk = ? and school_code = ?", nuptk, schoolCode).First(&t).Error
	return &t, err
}

// Create creates a new teacher
func (r *TeacherRepository) Create(t *teacher.Teacher) error {
	return r.db.Create(t).Error
}

// Update updates a teacher
func (r *TeacherRepository) Update(t *teacher.Teacher) error {
	return r.db.Save(t).Error
}

// Delete deletes a teacher
func (r *TeacherRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&teacher.Teacher{}).Error
}

// ----- Teacher Class Subject methods -----

// FindClassSubjectsByTeacher finds all class subjects for a teacher
func (r *TeacherRepository) FindClassSubjectsByTeacher(teacherID uint) ([]teacher.TeacherClassSubject, error) {
	var subjects []teacher.TeacherClassSubject
	err := r.db.Where("teacher_id = ?", teacherID).
		Preload("Subject").
		Preload("Class").
		Find(&subjects).Error
	return subjects, err
}

// FindClassSubjectByID finds a class subject by ID
func (r *TeacherRepository) FindClassSubjectByID(id uint) (*teacher.TeacherClassSubject, error) {
	var subject teacher.TeacherClassSubject
	err := r.db.Where("id = ?", id).First(&subject).Error
	return &subject, err
}

// FindClassSubjectDuplicate checks for duplicate class subject assignment
func (r *TeacherRepository) FindClassSubjectDuplicate(classID uint, subjectCode string, teacherID uint) (*teacher.TeacherClassSubject, error) {
	var subject teacher.TeacherClassSubject
	err := r.db.Where("class_id = ? AND subject_code = ? AND teacher_id = ?", classID, subjectCode, teacherID).First(&subject).Error
	return &subject, err
}

// CreateClassSubjects creates multiple class subject assignments
func (r *TeacherRepository) CreateClassSubjects(subjects []teacher.TeacherClassSubject) error {
	return r.db.Create(&subjects).Error
}

// DeleteClassSubject deletes a class subject assignment
func (r *TeacherRepository) DeleteClassSubject(id uint) error {
	return r.db.Where("id = ?", id).Delete(&teacher.TeacherClassSubject{}).Error
}
