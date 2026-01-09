package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/dapodik/entity"
	"gorm.io/gorm"
)

type DapodikRepository struct {
	db *gorm.DB
}

// NewDapodikRepository creates repository with injected DB
func NewDapodikRepository(db *gorm.DB) *DapodikRepository {
	return &DapodikRepository{db: db}
}

// DB returns the database connection for complex queries
func (r *DapodikRepository) DB() *gorm.DB {
	return r.db
}

// ----- Config methods -----

// FindConfigBySchoolCode finds dapodik config by school code
func (r *DapodikRepository) FindConfigBySchoolCode(schoolCode string) (*entity.DapodikConfig, error) {
	var config entity.DapodikConfig
	err := r.db.Where("school_code = ?", schoolCode).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateConfig creates a new dapodik config
func (r *DapodikRepository) CreateConfig(config *entity.DapodikConfig) error {
	return r.db.Create(config).Error
}

// UpdateConfig updates dapodik config
func (r *DapodikRepository) UpdateConfig(config *entity.DapodikConfig) error {
	return r.db.Save(config).Error
}

// ----- History methods -----

// FindRunningSync finds a running sync for school code
func (r *DapodikRepository) FindRunningSync(schoolCode string) (*entity.DapodikSyncHistory, error) {
	var history entity.DapodikSyncHistory
	err := r.db.Where("school_code = ? AND status IN ?", schoolCode, []string{
		string(entity.SyncStatusPending),
		string(entity.SyncStatusInProgress),
	}).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// CreateHistory creates a new sync history
func (r *DapodikRepository) CreateHistory(history *entity.DapodikSyncHistory) error {
	return r.db.Create(history).Error
}

// UpdateHistory updates sync history
func (r *DapodikRepository) UpdateHistory(history *entity.DapodikSyncHistory) error {
	return r.db.Save(history).Error
}

// FindHistoryByIDAndSchool finds sync history by ID and school code
func (r *DapodikRepository) FindHistoryByIDAndSchool(id uint, schoolCode string) (*entity.DapodikSyncHistory, error) {
	var history entity.DapodikSyncHistory
	err := r.db.Where("id = ? AND school_code = ?", id, schoolCode).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// CountBySchool counts total sync history by school code
func (r *DapodikRepository) CountBySchool(schoolCode string) int64 {
	var count int64
	r.db.Model(&entity.DapodikSyncHistory{}).Where("school_code = ?", schoolCode).Count(&count)
	return count
}

// CountBySchoolAndStatus counts sync history by school code and status
func (r *DapodikRepository) CountBySchoolAndStatus(schoolCode string, status entity.SyncStatus) int64 {
	var count int64
	r.db.Model(&entity.DapodikSyncHistory{}).Where("school_code = ? AND status = ?", schoolCode, status).Count(&count)
	return count
}

// FindLastSyncBySchool finds the last sync by school code
func (r *DapodikRepository) FindLastSyncBySchool(schoolCode string) (*entity.DapodikSyncHistory, error) {
	var history entity.DapodikSyncHistory
	err := r.db.Where("school_code = ?", schoolCode).Order("created_at DESC").First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}
