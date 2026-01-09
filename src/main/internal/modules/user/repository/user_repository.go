package repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/view"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates repository with injected DB
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// DB returns the database connection
func (r *UserRepository) DB() *gorm.DB {
	return r.db
}

// FindByID finds user by ID
func (r *UserRepository) FindByID(id uint) (*user.User, error) {
	var u user.User
	err := r.db.Where("id = ?", id).First(&u).Error
	return &u, err
}

// FindByUsername finds user by username
func (r *UserRepository) FindByUsername(username string) (*user.User, error) {
	var u user.User
	err := r.db.Where("username = ?", username).First(&u).Error
	return &u, err
}

// FindViewByID finds user view by ID
func (r *UserRepository) FindViewByID(id uint) (*view.VUser, error) {
	var u view.VUser
	err := r.db.Where("id = ?", id).First(&u).Error
	return &u, err
}

// Create creates a new user
func (r *UserRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

// Update updates a user
func (r *UserRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

// AllRoles returns all roles
func (r *UserRepository) AllRoles() ([]user.Role, error) {
	var roles []user.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

// FindRoleByCode finds role by code
func (r *UserRepository) FindRoleByCode(code string) (*user.Role, error) {
	var role user.Role
	err := r.db.Where("code = ?", code).First(&role).Error
	return &role, err
}
