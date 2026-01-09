package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"gorm.io/gorm"
)

type ClassRepository struct {
	db *gorm.DB
}

// NewClassRepository creates repository with injected DB
func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

// DB returns the database connection for complex queries
func (r *ClassRepository) DB() *gorm.DB {
	return r.db
}

// SchoolScope returns a scoped query for multi-tenant isolation
func (r *ClassRepository) SchoolScope(schoolCode string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("school_code = ?", schoolCode)
	}
}

// ----- Class methods -----

// FindByID finds class by ID
func (r *ClassRepository) FindByID(id uint) (*school.Class, error) {
	var class school.Class
	err := r.db.Where("id = ?", id).Preload("DetailClassCode").First(&class).Error
	return &class, err
}

// FindByIDScoped finds class by ID with school scope
func (r *ClassRepository) FindByIDScoped(schoolCode string, id uint) (*school.Class, error) {
	var class school.Class
	err := r.db.Scopes(r.SchoolScope(schoolCode)).
		Where("id = ?", id).Preload("DetailClassCode").First(&class).Error
	return &class, err
}

// FindByCodeScoped finds class by code with school scope
func (r *ClassRepository) FindByCodeScoped(schoolCode, classCode string) (*school.Class, error) {
	var class school.Class
	err := r.db.Scopes(r.SchoolScope(schoolCode)).
		Where("class_code = ?", classCode).First(&class).Error
	return &class, err
}

// Create creates a new class
func (r *ClassRepository) Create(class *school.Class) error {
	return r.db.Create(class).Error
}

// Update updates a class
func (r *ClassRepository) Update(class *school.Class) error {
	return r.db.Save(class).Error
}

// DeleteScoped deletes class by ID with school scope
func (r *ClassRepository) DeleteScoped(schoolCode string, id uint) error {
	return r.db.Scopes(r.SchoolScope(schoolCode)).
		Where("id = ?", id).Delete(&school.Class{}).Error
}

// ----- Teacher Class Subject methods -----

// FindTeacherClassSubjects finds all class subjects for a teacher
func (r *ClassRepository) FindTeacherClassSubjects(teacherID uint) ([]teacher.TeacherClassSubject, error) {
	var subjects []teacher.TeacherClassSubject
	err := r.db.Where("teacher_id = ?", teacherID).Find(&subjects).Error
	return subjects, err
}

// ----- Class Member methods -----

// FindMembersByClassID finds all students in a class
func (r *ClassRepository) FindMembersByClassID(classID uint) ([]view.VStudent, error) {
	var members []view.VStudent
	err := r.db.Where("class_id = ?", classID).Find(&members).Error
	return members, err
}

// FindStudentClassMembership checks if student is already in class
func (r *ClassRepository) FindStudentClassMembership(studentID, classID uint) (*student.StudentClass, error) {
	var member student.StudentClass
	err := r.db.Where("student_id = ? AND class_id = ?", studentID, classID).First(&member).Error
	return &member, err
}

// CreateClassMembers adds students to a class
func (r *ClassRepository) CreateClassMembers(members []student.StudentClass) error {
	return r.db.Create(&members).Error
}

// DeleteClassMemberByID removes a class member by ID
func (r *ClassRepository) DeleteClassMemberByID(id uint) error {
	return r.db.Where("id = ?", id).Delete(&student.StudentClass{}).Error
}

// DeleteClassMembersByClassID removes all class members by class ID
func (r *ClassRepository) DeleteClassMembersByClassID(classID uint) error {
	return r.db.Where("class_id = ?", classID).Delete(&student.StudentClass{}).Error
}

// DeleteClassMember removes a specific student from a class
func (r *ClassRepository) DeleteClassMember(studentID, classID uint) error {
	return r.db.Where("student_id = ? AND class_id = ?", studentID, classID).Delete(&student.StudentClass{}).Error
}
