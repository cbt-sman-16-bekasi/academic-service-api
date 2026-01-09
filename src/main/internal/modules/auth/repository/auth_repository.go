package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/teacher"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates repository with injected DB
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// DB returns the database connection
func (r *AuthRepository) DB() *gorm.DB {
	return r.db
}

// FindUserByUsername finds user by username with role preload
func (r *AuthRepository) FindUserByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.db.Where("username = ?", username).Preload("RoleUser").First(&u).Error
	return &u, err
}

// FindTeacherByUserID finds teacher by user ID
func (r *AuthRepository) FindTeacherByUserID(userID uint) (*teacher.Teacher, error) {
	var t teacher.Teacher
	err := r.db.Where("user_id = ?", userID).First(&t).Error
	return &t, err
}

// UpdateUser updates user data
func (r *AuthRepository) UpdateUser(u *user.User) error {
	return r.db.Save(u).Error
}

// UpdateTeacher updates teacher data
func (r *AuthRepository) UpdateTeacher(t *teacher.Teacher) error {
	return r.db.Save(t).Error
}
