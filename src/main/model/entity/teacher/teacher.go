package teacher

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"gorm.io/gorm"
)

type Teacher struct {
	UserId     uint       `gorm:"unique" json:"user_id"`
	DetailUser *user.User `json:"detail_user" gorm:"foreignKey:UserId;references:ID"`
	Name       string     `json:"name"`
	Nuptk      string     `json:"nuptk"`
	Gender     string     `json:"gender"`
	gorm.Model
}

func (t *Teacher) TableName() string {
	return TableNameTeacher
}
