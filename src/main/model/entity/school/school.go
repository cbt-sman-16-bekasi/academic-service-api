package school

import "gorm.io/gorm"

const TableNameSchool = "school_service.m_schools"

type School struct {
	gorm.Model
	SchoolCode string `gorm:"unique;not null" json:"school_code"`
	SchoolName string `gorm:"not null" json:"school_name"`
	Logo       string `gorm:"null" json:"logo"`
	Address    string `gorm:"null" json:"address"`
	Email      string `gorm:"null" json:"email"`
	Phone      string `gorm:"null" json:"phone"`
}

func (s *School) TableName() string {
	return TableNameSchool
}
