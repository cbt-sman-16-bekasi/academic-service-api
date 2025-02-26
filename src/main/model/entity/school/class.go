package school

import "gorm.io/gorm"

const (
	TableNameClass     = "school_service.m_class"
	TableNameClassCode = "school_service.m_class_code"
)

type Class struct {
	gorm.Model
	ClassCode string `json:"classCode"`
	ClassName string `json:"className"`
}

func (c *Class) TableName() string {
	return TableNameClass
}

type ClassCode struct {
	gorm.Model
	Code string `gorm:"unique" json:"code"`
	Name string `json:"name"`
}

func (c *ClassCode) TableName() string {
	return TableNameClassCode
}
