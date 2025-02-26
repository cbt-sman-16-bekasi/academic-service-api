package user

import "gorm.io/gorm"

type AccessRoleManagement struct {
	gorm.Model
	RoleCode         string           `gorm:"column:role_code" json:"roleCode"`
	PermissionCode   string           `gorm:"column:permission_code" json:"permissionCode"`
	PermissionAccess PermissionAccess `gorm:"foreignKey:PermissionCode" json:"permissionAccess"`
}

func (a *AccessRoleManagement) TableName() string {
	return TableNameAccessRoleManagement
}

type PermissionAccess struct {
	gorm.Model
	PermissionCode string                 `gorm:"column:permission_code" json:"permissionCode"`
	Name           string                 `gorm:"column:name" json:"name"`
	PermissionItem []PermissionItemAccess `gorm:"foreignKey:PermissionCode;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"permissionItem"`
}

func (p *PermissionAccess) TableName() string {
	return TableNamePermissionAccess
}

type PermissionItemAccess struct {
	gorm.Model
	PermissionCode string `gorm:"column:permission_code" json:"permissionCode"`
	MenuCode       string `gorm:"column:menu_code" json:"menuCode"`
	Menu           Menu   `gorm:"foreignKey:MenuCode" json:"menu"`
}

func (p *PermissionItemAccess) TableName() string {
	return TableNamePermissionItemAccess
}
