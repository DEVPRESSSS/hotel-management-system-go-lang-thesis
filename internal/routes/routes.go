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
		//Default route of html file while loading

		defaultRoute.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//This routes are temporary, it will be move later if there is authentication
		defaultRoute.POST("/userslist", server.CreateUser)
		defaultRoute.GET("/userslist", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "users.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		defaultRoute.GET("/roles", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "role.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all app users
		defaultRoute.GET("/api/users", server.GetAllUsers)
		//Fetch all roles
		defaultRoute.GET("/api/roles", server.GetRoles)

	}

}
