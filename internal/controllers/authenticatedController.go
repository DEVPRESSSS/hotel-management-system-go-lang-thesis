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
