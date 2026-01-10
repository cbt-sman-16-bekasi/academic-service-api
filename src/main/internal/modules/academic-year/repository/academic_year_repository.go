package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/academic-year/entity"
	"gorm.io/gorm"
)

// AcademicYearRepository handles database operations for academic year
type AcademicYearRepository struct {
	db *gorm.DB
}

// NewAcademicYearRepository creates a new repository instance
func NewAcademicYearRepository(db *gorm.DB) *AcademicYearRepository {
	return &AcademicYearRepository{db: db}
}

// DB returns the underlying database connection
func (r *AcademicYearRepository) DB() *gorm.DB {
	return r.db
}

// ===== Academic Year CRUD =====

// Create creates a new academic year
func (r *AcademicYearRepository) Create(academicYear *entity.AcademicYear) error {
	return r.db.Create(academicYear).Error
}

// Update updates an existing academic year
func (r *AcademicYearRepository) Update(academicYear *entity.AcademicYear) error {
	return r.db.Save(academicYear).Error
}

// Delete soft deletes an academic year
func (r *AcademicYearRepository) Delete(id uint) error {
	return r.db.Delete(&entity.AcademicYear{}, id).Error
}

// FindByID finds an academic year by ID
func (r *AcademicYearRepository) FindByID(id uint) (*entity.AcademicYear, error) {
	var academicYear entity.AcademicYear
	err := r.db.First(&academicYear, id).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

// FindBySemesterID finds an academic year by semester_id
func (r *AcademicYearRepository) FindBySemesterID(schoolCode, semesterID string) (*entity.AcademicYear, error) {
	var academicYear entity.AcademicYear
	err := r.db.Where("school_code = ? AND semester_id = ?", schoolCode, semesterID).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

// FindBySchoolCode finds all academic years for a school
func (r *AcademicYearRepository) FindBySchoolCode(schoolCode string) ([]entity.AcademicYear, error) {
	var academicYears []entity.AcademicYear
	err := r.db.Where("school_code = ?", schoolCode).Order("year_end DESC, semester DESC").Find(&academicYears).Error
	if err != nil {
		return nil, err
	}
	return academicYears, nil
}

// FindActive finds the active academic year for a school
func (r *AcademicYearRepository) FindActive(schoolCode string) (*entity.AcademicYear, error) {
	var academicYear entity.AcademicYear
	err := r.db.Where("school_code = ? AND is_active = ?", schoolCode, true).First(&academicYear).Error
	if err != nil {
		return nil, err
	}
	return &academicYear, nil
}

// SetActive sets an academic year as active and deactivates others
func (r *AcademicYearRepository) SetActive(schoolCode string, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Deactivate all other academic years for this school
		if err := tx.Model(&entity.AcademicYear{}).
			Where("school_code = ?", schoolCode).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// Activate the selected one
		if err := tx.Model(&entity.AcademicYear{}).
			Where("id = ?", id).
			Update("is_active", true).Error; err != nil {
			return err
		}

		return nil
	})
}

// ExistsBySemesterID checks if an academic year exists
func (r *AcademicYearRepository) ExistsBySemesterID(schoolCode, semesterID string) bool {
	var count int64
	r.db.Model(&entity.AcademicYear{}).
		Where("school_code = ? AND semester_id = ?", schoolCode, semesterID).
		Count(&count)
	return count > 0
}

// Count counts academic years for a school
func (r *AcademicYearRepository) Count(schoolCode string) int64 {
	var count int64
	r.db.Model(&entity.AcademicYear{}).Where("school_code = ?", schoolCode).Count(&count)
	return count
}

// ===== Student Class History CRUD =====

// CreateHistory creates a new student class history record
func (r *AcademicYearRepository) CreateHistory(history *entity.StudentClassHistory) error {
	return r.db.Create(history).Error
}

// UpdateHistory updates an existing history record
func (r *AcademicYearRepository) UpdateHistory(history *entity.StudentClassHistory) error {
	return r.db.Save(history).Error
}

// FindHistoryByID finds a history record by ID
func (r *AcademicYearRepository) FindHistoryByID(id uint) (*entity.StudentClassHistory, error) {
	var history entity.StudentClassHistory
	err := r.db.First(&history, id).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// FindHistoryByStudent finds all history records for a student
func (r *AcademicYearRepository) FindHistoryByStudent(schoolCode string, studentID uint) ([]entity.StudentClassHistory, error) {
	var histories []entity.StudentClassHistory
	err := r.db.Where("school_code = ? AND student_id = ?", schoolCode, studentID).
		Order("joined_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// FindHistoryByStudentAndSemester finds history for a student in a specific semester
func (r *AcademicYearRepository) FindHistoryByStudentAndSemester(schoolCode string, studentID uint, semesterID string) (*entity.StudentClassHistory, error) {
	var history entity.StudentClassHistory
	err := r.db.Where("school_code = ? AND student_id = ? AND semester_id = ?", schoolCode, studentID, semesterID).
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// FindHistoryByClass finds all history records for a class
func (r *AcademicYearRepository) FindHistoryByClass(schoolCode string, classID uint) ([]entity.StudentClassHistory, error) {
	var histories []entity.StudentClassHistory
	err := r.db.Where("school_code = ? AND class_id = ?", schoolCode, classID).
		Order("joined_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// FindHistoryBySemester finds all history records for a semester
func (r *AcademicYearRepository) FindHistoryBySemester(schoolCode, semesterID string) ([]entity.StudentClassHistory, error) {
	var histories []entity.StudentClassHistory
	err := r.db.Where("school_code = ? AND semester_id = ?", schoolCode, semesterID).
		Order("class_id, student_id").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// CountStudentsByClassAndSemester counts students in a class for a semester
func (r *AcademicYearRepository) CountStudentsByClassAndSemester(schoolCode string, classID uint, semesterID string) int64 {
	var count int64
	r.db.Model(&entity.StudentClassHistory{}).
		Where("school_code = ? AND class_id = ? AND semester_id = ? AND status = ?", schoolCode, classID, semesterID, "active").
		Count(&count)
	return count
}

// CountStudentsBySemester counts all students in a semester
func (r *AcademicYearRepository) CountStudentsBySemester(schoolCode, semesterID string) int64 {
	var count int64
	r.db.Model(&entity.StudentClassHistory{}).
		Where("school_code = ? AND semester_id = ? AND status = ?", schoolCode, semesterID, "active").
		Count(&count)
	return count
}

// FindOrCreateHistory finds or creates a student class history
func (r *AcademicYearRepository) FindOrCreateHistory(history *entity.StudentClassHistory) (*entity.StudentClassHistory, error) {
	existing, err := r.FindHistoryByStudentAndSemester(history.SchoolCode, history.StudentID, history.SemesterID)
	if err == nil {
		return existing, nil
	}

	if err := r.CreateHistory(history); err != nil {
		return nil, err
	}
	return history, nil
}

// BulkCreateHistory creates multiple history records
func (r *AcademicYearRepository) BulkCreateHistory(histories []entity.StudentClassHistory) error {
	if len(histories) == 0 {
		return nil
	}
	return r.db.CreateInBatches(histories, 100).Error
}
