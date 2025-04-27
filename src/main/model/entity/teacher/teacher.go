package teacher

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"gorm.io/gorm"
)

type Teacher struct {
	UserId       uint                  `gorm:"unique" json:"user_id"`
	DetailUser   *user.User            `json:"detail_user" gorm:"foreignKey:UserId;references:ID"`
	Name         string                `json:"name"`
	Nuptk        string                `json:"nuptk"`
	Gender       string                `json:"gender"`
	ClassSubject []TeacherClassSubject `json:"teacherClassSubject"`
	gorm.Model
}

func (t *Teacher) TableName() string {
	return TableNameTeacher
}

type TeacherClassSubject struct {
	TeacherId   uint               `json:"teacherId"`
	Teacher     Teacher            `json:"teacher" gorm:"foreignKey:TeacherId;references:ID"`
	SubjectCode string             `json:"subjectCode"`
	Subject     curriculum.Subject `json:"subject" gorm:"foreignKey:SubjectCode;references:Code"`
	ClassId     uint               `json:"classId"`
	Class       school.Class       `json:"class" gorm:"foreignKey:ClassId;references:ID"`
	gorm.Model
}

func (t *TeacherClassSubject) TableName() string {
	return TableNameTeacherClassSubject
}
