package student

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	SchoolCode string    `gorm:"type:varchar(50);index" json:"school_code"`
	UserId     uint      `json:"user_id"`
	Nisn       string    `json:"nisn"`
	DetailUser user.User `json:"detail_user" gorm:"foreignKey:UserId;references:ID"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
}

func (s *Student) TableName() string {
	return TableNameStudent
}
