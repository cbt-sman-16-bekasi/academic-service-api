package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	SchoolCode string `gorm:"not null" json:"SchoolCode"`
	Username   string `gorm:"unique" json:"username"`
	Password   string `gorm:"password" json:"-"`
	Salt       string `gorm:"salt_password" json:"-"`
	Role       uint   `gorm:"default:1" json:"-"`
	Status     uint   `gorm:"default:1" json:"status"`
	Name       string `json:"name"`
	RoleUser   Role   `gorm:"foreignKey:Role;references:ID" json:"role"`
}

func (u *User) TableName() string {
	return TableNameUser
}

type Role struct {
	gorm.Model
	SchoolCode string `gorm:"not null" json:"SchoolCode"`
	Code       string `gorm:"unique" json:"code"`
	Name       string `json:"name"`
}

func (r *Role) TableName() string {
	return TableNameRole
}
