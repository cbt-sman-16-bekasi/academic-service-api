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
}

func (s *StudentClass) TableName() string {
	return TableNameStudentClass
}
