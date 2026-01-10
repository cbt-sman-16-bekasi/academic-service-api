package student

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"gorm.io/gorm"
)

type StudentClass struct {
	gorm.Model
	SchoolCode    string       `gorm:"type:varchar(50);index" json:"school_code"`
	StudentId     uint         `json:"student_id"`
	DetailStudent Student      `json:"detail_student" gorm:"foreignKey:StudentId;references:ID"`
	ClassId       uint         `json:"class_id"`
	DetailClass   school.Class `json:"detail_class" gorm:"foreignKey:ClassId;references:ID"`
	// Semester tracking for Dapodik integration
	SemesterID          string `gorm:"type:varchar(10);index" json:"semester_id"`           // Dapodik format: YYYYS
	DapodikAnggotaID    string `gorm:"type:varchar(100);null" json:"dapodik_anggota_id"`    // anggota_rombel_id from Dapodik
	DapodikPesertaID    string `gorm:"type:varchar(100);null" json:"dapodik_peserta_id"`    // peserta_didik_id from Dapodik
	DapodikRegistrasiID string `gorm:"type:varchar(100);null" json:"dapodik_registrasi_id"` // registrasi_id from Dapodik
}

func (s *StudentClass) TableName() string {
	return TableNameStudentClass
}
