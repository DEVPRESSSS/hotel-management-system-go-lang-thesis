package routes

import (
	"HMS-GO/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(db *gorm.DB, router *gin.Engine) {

	// 404 handler
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{})
	})

	// Default routes

	defaultRoute := router.Group("/")
	{
		server := controllers.NewServer(db)

		// defaultRoute.GET("/", func(ctx *gin.Context) {
		// 	ctx.HTML(http.StatusOK, "defaultView/index.html", gin.H{
		// 		"title": "HoTel Management System",
		// 	})
		// })

		defaultRoute.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		defaultRoute.POST("/userslist", server.CreateUser)
		defaultRoute.GET("/userslist", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "users.html", gin.H{
				"title": "Hotel Management System",
			})
		})

	}

}
