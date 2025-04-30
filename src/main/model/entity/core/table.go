package core

type AuditUser struct {
	CreatedBy uint `json:"created_by"`
	//CreatedUser  user.User `json:"detail_user" gorm:"foreignKey:CreatedBy"`
	ModifiedBy uint `json:"modified_by"`
	//ModifiedUser user.User `json:"modified_user" gorm:"foreignKey:ModifiedBy"`
}
