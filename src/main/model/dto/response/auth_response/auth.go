package auth_response

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/cbt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
)

type AuthResponse struct {
	Token  string      `json:"token"`
	Exp    int64       `json:"exp"`
	User   *user.User  `json:"user"`
	Detail interface{} `json:"detail"`
}

type AuthResponseCBT struct {
	Token       string                        `json:"token,omitempty"`
	Exp         int64                         `json:"exp,omitempty"`
	User        *student.StudentClass         `json:"user,omitempty"`
	Exam        *school.Exam                  `json:"exam,omitempty"`
	ExamSession *school.ExamSession           `json:"exam_session,omitempty"`
	ExamTaken   *cbt.StudentHistoryTaken      `json:"exam_taken,omitempty"`
	ExamActive  []view.ExamSessionActiveToday `json:"exam_active,omitempty"`
}
