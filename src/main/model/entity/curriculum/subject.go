package curriculum

import "gorm.io/gorm"

const (
	TableNameSubject           = "curriculum_service.m_subject"
	TableNameCurriculumSubject = "curriculum_service.m_curriculum_subject"
)

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

type CurriculumSubject struct {
	gorm.Model
	SchoolCode   string `gorm:"type:varchar(50);index" json:"school_code"`
	CurriculumId uint   `json:"curriculum_id"`
	SubjectId    uint   `json:"subject_id"`
}

func (c *CurriculumSubject) TableName() string {
	return TableNameCurriculumSubject
}
