package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/curriculum/entity"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

// NewSubjectRepository creates repository with injected DB
func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

// DB returns the database connection for complex queries
func (r *SubjectRepository) DB() *gorm.DB {
	return r.db
}

// FindByID finds subject by ID
func (r *SubjectRepository) FindByID(id uint) (*entity.Subject, error) {
	var subject entity.Subject
	err := r.db.Where("id = ?", id).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// FindByCode finds subject by code
func (r *SubjectRepository) FindByCode(code string) (*entity.Subject, error) {
	var subject entity.Subject
	err := r.db.Where("code = ?", code).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// FindViewByID finds subject view by ID
func (r *SubjectRepository) FindViewByID(id uint) (*entity.VSubject, error) {
	var subject entity.VSubject
	err := r.db.Where("id = ?", id).First(&subject).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// Create creates a new subject
func (r *SubjectRepository) Create(subject *entity.Subject) error {
	return r.db.Create(subject).Error
}

// Update updates a subject
func (r *SubjectRepository) Update(subject *entity.Subject) error {
	return r.db.Save(subject).Error
}

// Delete soft deletes a subject
func (r *SubjectRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&entity.Subject{}).Error
}

// Exists checks if subject with code exists
func (r *SubjectRepository) Exists(schoolCode, code string) bool {
	var count int64
	r.db.Model(&entity.Subject{}).Where("code = ? and school_code = ?", code, schoolCode).Count(&count)
	return count > 0
}

// ExistsExcluding checks if subject with code exists excluding given ID
func (r *SubjectRepository) ExistsExcluding(code string, excludeID uint) bool {
	var count int64
	r.db.Model(&entity.Subject{}).Where("code = ? AND id != ?", code, excludeID).Count(&count)
	return count > 0
}
