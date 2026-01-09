package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	TableNameDapodikConfig      = "school_service.m_dapodik_config"
	TableNameDapodikSyncHistory = "school_service.m_dapodik_sync_history"
)

// SyncType menentukan tipe sinkronisasi
type SyncType string

const (
	SyncTypeUpdate    SyncType = "UPDATE"     // Update & Tambah Data Baru (direkomendasikan)
	SyncTypeFullReset SyncType = "FULL_RESET" // Reset & Import Ulang Semua Data
)

// SyncStatus menentukan status sinkronisasi
type SyncStatus string

const (
	SyncStatusPending    SyncStatus = "PENDING"
	SyncStatusInProgress SyncStatus = "IN_PROGRESS"
	SyncStatusSuccess    SyncStatus = "SUCCESS"
	SyncStatusFailed     SyncStatus = "FAILED"
)

// DapodikConfig menyimpan konfigurasi koneksi ke API DAPODIK
type DapodikConfig struct {
	gorm.Model
	SchoolCode  string     `gorm:"type:varchar(20);not null;unique" json:"school_code"`
	DapodikURL  string     `gorm:"type:varchar(255);not null" json:"dapodik_url"`
	APIKey      string     `gorm:"type:varchar(500);not null" json:"api_key"`
	UseProxy    bool       `gorm:"default:false" json:"use_proxy"`              // Toggle: true = via proxy, false = direct
	ProxyURL    string     `gorm:"type:varchar(255);null" json:"proxy_url"`     // Optional: proxy server URL
	ProxyAPIKey string     `gorm:"type:varchar(255);null" json:"proxy_api_key"` // Optional: API key for proxy
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	LastSyncAt  *time.Time `gorm:"null" json:"last_sync_at"`
}

func (d *DapodikConfig) TableName() string {
	return TableNameDapodikConfig
}

// DapodikSyncHistory menyimpan riwayat sinkronisasi dengan DAPODIK
type DapodikSyncHistory struct {
	gorm.Model
	SchoolCode        string     `gorm:"type:varchar(20);not null;index" json:"school_code"`
	SyncType          SyncType   `gorm:"type:varchar(20);not null" json:"sync_type"`
	Status            SyncStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	StartedAt         *time.Time `gorm:"null" json:"started_at"`
	CompletedAt       *time.Time `gorm:"null" json:"completed_at"`
	TotalStudents     int        `gorm:"default:0" json:"total_students"`
	TotalClasses      int        `gorm:"default:0" json:"total_classes"`
	TotalTeachers     int        `gorm:"default:0" json:"total_teachers"`
	SyncedStudents    int        `gorm:"default:0" json:"synced_students"`
	SyncedClasses     int        `gorm:"default:0" json:"synced_classes"`
	SyncedTeachers    int        `gorm:"default:0" json:"synced_teachers"`
	ErrorMessage      string     `gorm:"type:text;null" json:"error_message"`
	TriggeredByUserId uint       `gorm:"null" json:"triggered_by_user_id"`
	TriggeredByName   string     `gorm:"type:varchar(100);null" json:"triggered_by_name"`
}

func (d *DapodikSyncHistory) TableName() string {
	return TableNameDapodikSyncHistory
}
