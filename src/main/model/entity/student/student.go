package student

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"gorm.io/gorm"
)

type Student struct {
	UserId     uint      `json:"user_id"`
	DetailUser user.User `json:"detail_user" gorm:"foreignKey:UserId;references:ID"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender"`
	gorm.Model
}

func (s *Student) TableName() string {
	return TableNameStudent
}
