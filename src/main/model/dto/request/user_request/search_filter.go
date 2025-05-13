package user_request

type UserSearchRequest struct {
	Name     string `form:"name" json:"name"`
	Username string `form:"username" json:"username"`
	Status   string `form:"status" json:"status"`
}
