package entity

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/modules/spp/entity"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/database"
)

func init() {
	database.MigrationRegister(
		// SPP Module
		&entity.Payment{},
		&entity.PaymentTransaction{},
		&entity.PaymentConfig{},
	//&user.User{},
	//&user.Role{},
	//&user.AccessRoleManagement{},
	//&user.PermissionAccess{},
	//&user.PermissionItemAccess{},
	//&user.Menu{},
	//&school.Class{},
	//&school.ClassSubject{},
	//&school.SchoolConfig{},
	//&school.SchoolLevel{},
	//&school.School{},
	//&curriculum.Curriculum{},
	//&curriculum.CurriculumSubject{},
	//&curriculum.Subject{},
	//&teacher.Teacher{},
	//&school.TypeExam{},
	//&school.Exam{},
	//&school.ExamMember{},
	//&school.ExamQuestion{},
	//&school.ExamAnswerOption{},
	//&school.ExamSession{},
	//&student.Student{},
	//&student.StudentClass{},
	//&school.TokenExamSession{},
	//&cbt.StudentAnswers{},
	//&cbt.StudentHistoryTaken{},
	//&school.MasterBankQuestion{},
	//&school.BankQuestion{},
	//&school.BankAnswerOption{},
	//&teacher.TeacherClassSubject{},
	//&school.TokenExamSession{},
	//&school.ExamSessionMember{},
	)
}
