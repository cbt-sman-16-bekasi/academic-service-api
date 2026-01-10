package school

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"gorm.io/gorm"
)

const (
	TableNameClass        = "school_service.m_class"
	TableNameClassCode    = "school_service.m_class_code"
	TableNameClassSubject = "school_service.m_class_subject"
)

type Class struct {
	gorm.Model
	SchoolCode      string    `gorm:"type:varchar(50);index" json:"school_code"`
	ClassCode       string    `json:"classCode"`
	DetailClassCode ClassCode `gorm:"foreignKey:ClassCode;references:Code"`
	ClassName       string    `json:"className"`
	// Semester tracking for Dapodik integration
	SemesterID      string `gorm:"type:varchar(10);index" json:"semester_id"`       // Dapodik format: YYYYS (e.g., "20251")
	DapodikRombelID string `gorm:"type:varchar(100);null" json:"dapodik_rombel_id"` // rombongan_belajar_id from Dapodik
	WaliKelasID     *uint  `gorm:"null" json:"wali_kelas_id"`                       // Teacher ID for homeroom teacher
	WaliKelasName   string `gorm:"type:varchar(100);null" json:"wali_kelas_name"`   // Cached name for display
}

func (c *Class) TableName() string {
	return TableNameClass
}

type ClassCode struct {
	gorm.Model
	SchoolCode  string  `gorm:"type:varchar(50);index" json:"school_code"`
	Code        string  `gorm:"unique" json:"code"`
	Name        string  `json:"name"`
	ClassMember []Class `json:"class_member" gorm:"foreignKey:ClassCode;references:Code"`
}

func (c *ClassCode) TableName() string {
	return TableNameClassCode
}

type ClassSubject struct {
	gorm.Model
	SchoolCode      string             `gorm:"type:varchar(50);index" json:"school_code"`
	SubjectCode     string             `json:"subjectCode"`
	DetailSubject   curriculum.Subject `gorm:"foreignKey:SubjectCode;references:Code"`
	ClassCode       string             `json:"classCode"`
	DetailClassCode ClassCode          `gorm:"foreignKey:ClassCode;references:Code"`
}

func (c *ClassSubject) TableName() string {
	return TableNameClassSubject
}
