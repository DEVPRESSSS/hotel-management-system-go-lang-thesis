package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderWithRole renders HTML template with role automatically injected
func RenderWithRole(ctx *gin.Context, templateName string, data gin.H) {
	roleInterface, _ := ctx.Get("role")
	role := ""
	if r, ok := roleInterface.(string); ok {
		role = r
	}

	// Ensure data map exists
	if data == nil {
		data = gin.H{}
	}

	// Add role to data
	data["role"] = role

	ctx.HTML(http.StatusOK, templateName, data)
}
