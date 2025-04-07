package entity

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/cbt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/curriculum"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/yon-module/yon-framework/database"
)

func init() {
	database.MigrationRegister(
		&user.User{},
		&user.Role{},
		&user.AccessRoleManagement{},
		&user.PermissionAccess{},
		&user.PermissionItemAccess{},
		&user.Menu{},
		&school.Class{},
		&school.ClassSubject{},
		&school.SchoolConfig{},
		&school.SchoolLevel{},
		&school.School{},
		&curriculum.Curriculum{},
		&curriculum.CurriculumSubject{},
		&curriculum.Subject{},
		&teacher.Teacher{},
		&school.TypeExam{},
		&school.Exam{},
		&school.ExamMember{},
		&school.ExamQuestion{},
		&school.ExamAnswerOption{},
		&school.ExamSession{},
		&student.Student{},
		&student.StudentClass{},
		&school.TokenExamSession{},
		&cbt.StudentAnswers{},
		&cbt.StudentHistoryTaken{},
		&school.MasterBankQuestion{},
		&school.BankQuestion{},
		&school.BankAnswerOption{},
	)
}
