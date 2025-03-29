package school

import "gorm.io/gorm"

const TableNameSchool = "school_service.m_schools"

type School struct {
	gorm.Model
	SchoolCode string `gorm:"unique;not null" json:"school_code"`
	SchoolName string `gorm:"not null" json:"school_name"`
	NSS        string `json:"nss"`
	NPSN       string `json:"npsn"`
	Logo       string `gorm:"null" json:"logo"`
	Banner     string `gorm:"null" json:"banner"`
	Address    string `gorm:"null" json:"address"`
	Email      string `gorm:"null" json:"email"`
	Phone      string `gorm:"null" json:"phone"`
}

func (s *School) TableName() string {
	return TableNameSchool
}
