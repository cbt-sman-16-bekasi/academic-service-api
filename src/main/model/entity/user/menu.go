package user

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	SchoolCode string `gorm:"not null" json:"SchoolCode"`
	MenuCode   string `gorm:"unique" json:"menuCode"`
	ParentCode string `json:"parentCode"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
}

func (m *Menu) TableName() string {
	return TableNameMenu
}
