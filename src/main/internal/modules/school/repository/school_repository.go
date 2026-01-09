package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"gorm.io/gorm"
)

type SchoolRepository struct {
	db *gorm.DB
}

// NewSchoolRepository creates repository with injected DB
func NewSchoolRepository(db *gorm.DB) *SchoolRepository {
	return &SchoolRepository{db: db}
}

// DB returns the database connection
func (r *SchoolRepository) DB() *gorm.DB {
	return r.db
}

// FindSchoolByCode finds school by school code
func (r *SchoolRepository) FindSchoolByCode(schoolCode string) (*school.School, error) {
	var s school.School
	err := r.db.Where("school_code = ?", schoolCode).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// UpdateSchool updates school data
func (r *SchoolRepository) UpdateSchool(s *school.School) error {
	return r.db.Save(s).Error
}

// GetAllClassCodes returns all class codes
func (r *SchoolRepository) GetAllClassCodes() []school.ClassCode {
	var classCodes []school.ClassCode
	r.db.Find(&classCodes)
	return classCodes
}

// AllClassCodeSchool returns all class codes
func (r *SchoolRepository) AllClassCodeSchool(schoolCode string) []school.ClassCode {
	var classCodes []school.ClassCode
	r.db.Where("school_code", schoolCode).Find(&classCodes)
	return classCodes
}

// GetAllSubjects returns all subjects
func (r *SchoolRepository) GetAllSubjects() []curriculum.Subject {
	var subjects []curriculum.Subject
	r.db.Find(&subjects)
	return subjects
}

// FindClassSubjectByID finds class-subject mapping by ID
func (r *SchoolRepository) FindClassSubjectByID(id uint) (*school.ClassSubject, error) {
	var cs school.ClassSubject
	err := r.db.Preload("DetailSubject").Preload("DetailClassCode").Where("id = ?", id).First(&cs).Error
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// FindClassSubjectByClassAndSubject finds existing mapping
func (r *SchoolRepository) FindClassSubjectByClassAndSubject(classCode, subjectCode string) (*school.ClassSubject, error) {
	var cs school.ClassSubject
	err := r.db.Where("class_code = ? AND subject_code = ?", classCode, subjectCode).First(&cs).Error
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// FindSubjectByCode finds subject by code
func (r *SchoolRepository) FindSubjectByCode(code string) (*curriculum.Subject, error) {
	var s curriculum.Subject
	err := r.db.Where("code = ?", code).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateClassSubject creates new class-subject mapping
func (r *SchoolRepository) CreateClassSubject(cs *school.ClassSubject) error {
	return r.db.Create(cs).Error
}

// UpdateClassSubject updates class-subject mapping
func (r *SchoolRepository) UpdateClassSubject(cs *school.ClassSubject) error {
	return r.db.Save(cs).Error
}

// DeleteClassSubject deletes class-subject mapping
func (r *SchoolRepository) DeleteClassSubject(id uint) error {
	return r.db.Delete(&school.ClassSubject{}, id).Error
}

// GetDashboardAdmin returns dashboard data for admin
func (r *SchoolRepository) GetDashboardAdmin() (*view.DashboardSummary, error) {
	var dashboard view.DashboardSummary
	err := r.db.First(&dashboard).Error
	return &dashboard, err
}

// GetDashboardTeacher returns dashboard data for teacher
func (r *SchoolRepository) GetDashboardTeacher(teacherID uint) (*view.DashboardTeacher, error) {
	var dashboard view.DashboardTeacher
	err := r.db.Where("teacher_id = ?", teacherID).First(&dashboard).Error
	return &dashboard, err
}

// FindSystemConfigByOrigin finds system config by origin
func (r *SchoolRepository) FindSystemConfigByOrigin(origin string) (*school.SystemConfig, error) {
	var config school.SystemConfig
	err := r.db.Where("origin = ?", origin).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
