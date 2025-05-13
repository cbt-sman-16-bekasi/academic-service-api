package school_repository

import (
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/user"
	"github.com/yon-module/yon-framework/database"
	"gorm.io/gorm"
)

type UserRepository struct {
	Database *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		Database: database.GetDB(),
	}
}

func (r *UserRepository) ReadUser(username string) (user *user.User) {
	r.Database.Where("username = ?", username).Preload("RoleUser").First(&user)
	if user.ID == 0 {
		user = nil
	}
	return
}

func (r *UserRepository) FindById(id uint) (user *user.User) {
	r.Database.Where("id = ? AND status = '1'", id).Preload("RoleUser").First(&user)
	if user.ID == 0 {
		user = nil
	}
	return
}

func (r *UserRepository) AllRole() []user.Role {
	var roles []user.Role
	r.Database.Find(&roles)
	return roles
}

func (r *UserRepository) ReadRole(code string) *user.Role {
	var role user.Role
	r.Database.Where("code = ?", code).First(&role)
	if role.ID == 0 {
		return nil
	}
	return &role
}
