package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"HMS-GO/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Login(ctx *gin.Context) {
	var user dto.LoginInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, role, err := s.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie(
		"token",
		token,
		3600,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"success": "Account found", "token": token, "role": role})

}

func (s *Server) AuthenticateUser(username, password string) (string, string, error) {
	var user models.User

	if err := s.Db.
		Preload("Role").
		Where("username = ?", username).
		Take(&user).Error; err != nil {
		return "", "", fmt.Errorf("Incorrect username or password")
	}

	if err := utils.VerifyPassword(user.Password, password); err != nil {
		return "", "", fmt.Errorf("Incorrect username or password")
	}

	//Load permissions
	var roleAccess []models.RoleAccess
	if err := s.Db.
		Preload("Access").
		Where("role_id = ?", user.RoleId).
		Find(&roleAccess).Error; err != nil {

		return "", "", fmt.Errorf("failed to load permissions!!!!")
	}

	token, err := utils.CreateToken(user, roleAccess)
	if err != nil {

		return "", "", fmt.Errorf("failed to generate token")

	}

	return token, user.Role.RoleName, nil
}


