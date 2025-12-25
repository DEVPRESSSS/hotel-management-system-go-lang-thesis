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
	server := controllers.NewServer(db)

	defaultRoute := router.Group("/")
	{

		//Default route of html file while loading
		defaultRoute.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//This routes are temporary, it will be move later if there is authentication
		defaultRoute.POST("/userslist", server.CreateUser)
		defaultRoute.PUT("/api/update/:userid", server.UpdateUser)
		defaultRoute.DELETE("/api/delete/:userid", server.DeleteUser)

		defaultRoute.GET("/userslist", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "users.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Render role.html file
		defaultRoute.POST("/api/create", server.CreateRole)
		defaultRoute.GET("/roles", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "role.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all app users
		defaultRoute.GET("/api/users", server.GetAllUsers)
		//Fetch selected user information
		defaultRoute.GET("/api/user/:userid", server.GetUser)

		//Fetch all roles
		defaultRoute.GET("/api/roles", server.GetRoles)
		//Fetch selected role information
		defaultRoute.GET("/api/roles/:roleid", server.GetRole)

	}
}
