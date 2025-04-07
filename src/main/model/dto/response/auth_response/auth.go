package auth_response

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/cbt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
)

type AuthResponse struct {
	Token string     `json:"token"`
	Exp   int64      `json:"exp"`
	User  *user.User `json:"user"`
}

type AuthResponseCBT struct {
	Token       string                   `json:"token"`
	Exp         int64                    `json:"exp"`
	User        *student.StudentClass    `json:"user"`
	Exam        *school.Exam             `json:"exam"`
	ExamSession *school.ExamSession      `json:"exam_session"`
	ExamTaken   *cbt.StudentHistoryTaken `json:"exam_taken"`
}
