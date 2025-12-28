package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Login(ctx *gin.Context) {

	var user models.LoginInput

	if err := ctx.ShouldBindJSON(&user); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err := s.AuthenticateUser(user.Username, user.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "Account found"})
}

func (s *Server) AuthenticateUser(username, password string) (string, error) {
	var user models.User

	if err := s.Db.
		Where("username = ?", username).
		Take(&user).Error; err != nil {
		return "", fmt.Errorf("account not found")
	}

	if err := utils.VerifyPassword(user.Password, password); err != nil {
		return "", fmt.Errorf("incorrect password")
	}

	return user.Username, nil
}
