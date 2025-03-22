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
	ClassCode       string    `json:"classCode"`
	DetailClassCode ClassCode `gorm:"foreignKey:ClassCode"`
	ClassName       string    `json:"className"`
}

func (c *Class) TableName() string {
	return TableNameClass
}

type ClassCode struct {
	gorm.Model
	Code string `gorm:"unique" json:"code"`
	Name string `json:"name"`
}

func (c *ClassCode) TableName() string {
	return TableNameClassCode
}

type ClassSubject struct {
	gorm.Model
	SubjectCode     string             `json:"subjectCode"`
	DetailSubject   curriculum.Subject `gorm:"foreignKey:SubjectCode"`
	ClassCode       string             `json:"classCode"`
	DetailClassCode ClassCode          `gorm:"foreignKey:ClassCode"`
}

func (c *ClassSubject) TableName() string {
	return TableNameClassSubject
}
