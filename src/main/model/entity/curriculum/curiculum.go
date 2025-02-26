package curriculum

import "gorm.io/gorm"

const TableNameCurriculum = "curriculum_service.m_curriculum"

type Curriculum struct {
	gorm.Model
	Code   string `gorm:"unique" json:"code"`
	Name   string `gorm:"unique" json:"name"`
	Status bool   `gorm:"default:false" json:"status"`
}

func (c *Curriculum) TableName() string {
	return TableNameCurriculum
}
