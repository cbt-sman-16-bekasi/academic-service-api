package entity

import "gorm.io/gorm"

const (
	TableNameSubject           = "curriculum_service.m_subject"
	TableNameCurriculumSubject = "curriculum_service.m_curriculum_subject"
	TableNameCurriculum        = "curriculum_service.m_curriculum"
)

// Subject represents a curriculum subject
type Subject struct {
	gorm.Model
	SchoolCode  string `gorm:"type:varchar(50);index" json:"school_code"`
	Code        string `gorm:"unique" json:"code"`
	Subject     string `json:"subject"`
	SubjectType string `json:"subject_type"`
	Description string `json:"description"`
}

func (s *Subject) TableName() string {
	return TableNameSubject
}

// Curriculum represents a curriculum
type Curriculum struct {
	gorm.Model
	Code   string `gorm:"unique" json:"code"`
	Name   string `gorm:"unique" json:"name"`
	Status bool   `gorm:"default:false" json:"status"`
}

func (c *Curriculum) TableName() string {
	return TableNameCurriculum
}

// CurriculumSubject represents curriculum-subject relation
type CurriculumSubject struct {
	gorm.Model
	SchoolCode   string `gorm:"type:varchar(50);index" json:"school_code"`
	CurriculumId uint   `json:"curriculum_id"`
	SubjectId    uint   `json:"subject_id"`
}

func (c *CurriculumSubject) TableName() string {
	return TableNameCurriculumSubject
}

// VSubject is a view for subject with class codes
type VSubject struct {
	Subject
	ClassCode string `json:"class_code" gorm:"column:class_code"`
}

func (v *VSubject) TableName() string {
	return "public.v_subject"
}

