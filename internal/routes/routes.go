package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(db *gorm.DB, router *gin.Engine) {

	// 404 handler
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "error.html", gin.H{})
	})

	// Default routes
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "defaultView/index.html", gin.H{
			"title": "HoTel Management System",
		})
	})
	
	
}
