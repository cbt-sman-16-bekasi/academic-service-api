package dto

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/cbt"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/student"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
)

// ----- Request DTOs -----

// AuthRequest for admin/teacher login
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CBTAuthRequest for student CBT login
type CBTAuthRequest struct {
	Username string `json:"username" binding:"required"`
}

// ChangePasswordRequest for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// ChangeProfileRequest for changing profile
type ChangeProfileRequest struct {
	FullName   string `json:"full_name"`
	Username   string `json:"username"`
	ProfileURL string `json:"profile_url"`
}

// CBTSelectedSession for selecting exam session
type CBTSelectedSession struct {
	SessionID string `json:"session_id" binding:"required"`
}

// ----- Response DTOs -----

// AuthResponse for admin/teacher login response
type AuthResponse struct {
	Token  string      `json:"token"`
	Exp    int64       `json:"exp"`
	User   *user.User  `json:"user"`
	Detail interface{} `json:"detail"`
}

// AuthResponseCBT for student CBT login response
type AuthResponseCBT struct {
	Token       string                        `json:"token,omitempty"`
	Exp         int64                         `json:"exp,omitempty"`
	User        *student.StudentClass         `json:"user,omitempty"`
	Exam        *school.Exam                  `json:"exam,omitempty"`
	ExamSession *school.ExamSession           `json:"exam_session,omitempty"`
	ExamTaken   *cbt.StudentHistoryTaken      `json:"exam_taken,omitempty"`
	ExamActive  []view.ExamSessionActiveToday `json:"exam_active,omitempty"`
}


