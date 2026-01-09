package dto

// ----- Request DTOs -----

// UserUpdateRequest for create/update user
type UserUpdateRequest struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name"`
	Role     *int   `json:"role"`
	Status   int    `json:"status"`
}

// ----- Response DTOs -----

// RoleResponse for role data
type RoleResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// UserResponse for user data
type UserResponse struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Status     uint   `json:"status"`
	SchoolCode string `json:"school_code"`
}
