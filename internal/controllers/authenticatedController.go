package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

func (s *Server) GetUserId(ctx *gin.Context) string {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return ""
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return ""
	}
	return userIdStr
}
