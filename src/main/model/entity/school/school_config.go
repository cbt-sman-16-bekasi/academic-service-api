package school

import "gorm.io/gorm"

const (
	TableNameSchoolConfig      = "school_service.m_school_config"
	TableNameSchoolLevelConfig = "school_service.m_school_level_config"
)

type SchoolConfig struct {
	gorm.Model
	SchoolCode    string `gorm:"type:varchar(20);unique_index"`
	SchoolLevelId uint   `gorm:"not null" json:"schoolLevelId"`
}

func (c *SchoolConfig) TableName() string {
	return TableNameSchoolConfig
}

type SchoolLevel struct {
	gorm.Model
	Name string `gorm:"not null" json:"name"`
}

func (s *SchoolLevel) TableName() string {
	return TableNameSchoolLevelConfig
}
