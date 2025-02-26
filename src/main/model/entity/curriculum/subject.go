package curriculum

import "gorm.io/gorm"

const (
	TableNameSubject           = "curriculum_service.m_subject"
	TableNameCurriculumSubject = "curriculum_service.m_curriculum_subject"
)

type Subject struct {
	gorm.Model
	Code        string `gorm:"unique" json:"code"`
	Subject     string `gorm:"unique" json:"subject"`
	SubjectType string `json:"subject_type"`
	Description string `json:"description"`
}

func (s *Subject) TableName() string {
	return TableNameSubject
}

type CurriculumSubject struct {
	gorm.Model
	CurriculumId uint `json:"curriculum_id"`
	SubjectId    uint `json:"subject_id"`
}

func (c *CurriculumSubject) TableName() string {
	return TableNameCurriculumSubject
}
