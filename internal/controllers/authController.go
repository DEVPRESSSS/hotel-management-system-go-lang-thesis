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
	token, err := s.AuthenticateUser(user.Username, user.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		fmt.Print(err)
		return
	}

	//Set the cookie here
	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
	//Show status code OK and a success message
	ctx.JSON(http.StatusOK, gin.H{"success": "Account found", "token": token})
}

func (s *Server) AuthenticateUser(username, password string) (string, error) {
	var user models.User

	if err := s.Db.
		Preload("Role").
		Where("username = ?", username).
		Take(&user).Error; err != nil {
		return "", fmt.Errorf("Incorrect username or password")
	}

	if err := utils.VerifyPassword(user.Password, password); err != nil {
		return "", fmt.Errorf("Incorrect username or password")
	}

	token, err := utils.CreateToken(user)
	if err != nil {

		return "", fmt.Errorf("failed to generate token")

	}

	return token, nil
}
