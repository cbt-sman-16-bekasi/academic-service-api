package user_request

type UserUpdateRequest struct {
	Name     string `form:"name" json:"name"`
	Username string `form:"username" json:"username"`
	Status   int    `form:"status" json:"status"`
}
