package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/models/dto"
	"HMS-GO/internal/utils"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
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

func (s *Server) ForgotPassword(ctx *gin.Context) {
	var payload *models.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"
	var user models.User
	result := s.Db.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Email does not exist!"})
		return
	}

	// if !user.Verified {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
	// 	return
	// }

	// Generate Verification Code
	resetToken := randstr.String(20)
	passwordResetToken := utils.Encode(resetToken)
	user.PasswordResetToken = passwordResetToken
	user.PasswordResetAt = time.Now().Add(time.Minute * 15)
	s.Db.Save(&user)

	var firstName = user.FullName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// Get client origin from environment variable
	clientOrigin := os.Getenv("BASE_URL")

	// Send Email
	emailData := utils.EmailData{
		URL:       clientOrigin + "/resetpassword-form/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 15min)",
	}
	utils.SendEmail(&user, &emailData, "resetPassword.html")

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (s *Server) ResetPassword(ctx *gin.Context) {
	var payload *models.ResetPasswordInput
	resetToken := ctx.Params.ByName("resetToken")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if payload.Password != payload.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	passwordResetToken := utils.Encode(resetToken)

	var updatedUser models.User
	result := s.Db.First(&updatedUser, "password_reset_token = ? AND password_reset_at > ?", passwordResetToken, time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "The reset token is invalid or has expired"})
		return
	}

	// Set the plain password, then hash it
	updatedUser.Password = payload.Password
	if err := utils.HashPassword(&updatedUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Failed to hash password"})
		return
	}

	// Clear the reset token
	updatedUser.PasswordResetToken = ""
	s.Db.Save(&updatedUser)

	// Clear the auth cookie
	ctx.SetCookie("token", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password updated successfully"})
}
