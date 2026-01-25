package dto

type LoginInput struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}
