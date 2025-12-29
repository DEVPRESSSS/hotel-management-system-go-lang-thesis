package routes

import (
	"HMS-GO/internal/controllers"
	"HMS-GO/internal/middlewares"
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

		/*
			---------Authentication route-------
		*/
		defaultRoute.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		defaultRoute.POST("/api/auth", server.Login)

		/*
			END
		*/

	}

	authorize := router.Group("/hms")
	authorize.Use(middlewares.AuthMiddleware())
	{

		/*
			---------CRUD USER-------
		*/

		//Create user route
		authorize.POST("/userslist", server.CreateUser)

		//Update user route
		authorize.PUT("/api/update/:userid", server.UpdateUser)

		//Delete user route
		authorize.DELETE("/api/delete/:userid", server.DeleteUser)

		//Route for users list
		authorize.GET("/userslist", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "users.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		//Populate user.html using api/user API
		authorize.GET("/api/users", server.GetAllUsers)

		//Fetch selected user from API
		authorize.GET("/api/user/:userid", server.GetUser)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD ROLE-------
		*/

		//Create role route
		authorize.POST("/api/createrole", server.CreateRole)

		//Update role route
		authorize.PUT("/api/updaterole/:roleid", server.UpdateRole)

		//Delete role route
		authorize.DELETE("/api/deleterole/:roleid", server.DeleteRole)

		//Route for roles list
		authorize.GET("/roles", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "role.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all roles
		authorize.GET("/api/roles", server.GetRoles)

		//Fetch selected role information
		authorize.GET("/api/roles/:roleid", server.GetRole)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD FACILITY-------
		*/
		//Create facility route
		authorize.POST("/api/createfacility", server.CreateFacility)

		//Update facility route
		authorize.PUT("/api/updatefacility/:facilityid", server.UpdateFacility)

		//Delete facility route
		authorize.DELETE("/api/deletefacility/:facilityid", server.Deletefacility)

		//Routes for facility
		authorize.GET("/facility", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "facility.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all roles
		authorize.GET("/api/facility", server.GetFacilities)

		//Fetch selected role information
		authorize.GET("/api/facility/:facilityid", server.GetFacility)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD SERVICE-------
		*/
		//Create services route
		authorize.POST("/api/createservice", server.CreateService)

		//Update services route
		authorize.PUT("/api/updateservice/:serviceid", server.UpdateService)

		//Delete services route
		authorize.DELETE("/api/deleteservice/:serviceid", server.DeleteService)

		//Routes for services
		authorize.GET("/service", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "service.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all services
		authorize.GET("/api/services", server.GetServices)

		//Fetch selected services information
		authorize.GET("/api/service/:serviceid", server.GetService)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD SERVICE-------
		*/
		//Create room route
		authorize.POST("/api/createroom", server.CreateRoom)

		//Update room route
		authorize.PUT("/api/updateroom/:roomid", server.UpdateRoom)

		//Delete room route
		authorize.DELETE("/api/deleteroom/:roomid", server.DeleteRoom)

		//Routes for room
		authorize.GET("/rooms", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "room.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all room
		authorize.GET("/api/rooms", server.GetRooms)

		//Fetch selected room information
		authorize.GET("/api/room/:roomid", server.GetRoom)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD FLOOR-------
		*/
		// //Create room route
		// defaultRoute.POST("/api/createroom", server.CreateRoom)

		// //Update room route
		// defaultRoute.PUT("/api/updateroom/:roomid", server.UpdateRoom)

		// //Delete room route
		// defaultRoute.DELETE("/api/deleteroom/:roomid", server.DeleteRoom)

		//Routes for room
		// defaultRoute.GET("/rooms", func(ctx *gin.Context) {
		// 	ctx.HTML(http.StatusOK, "room.html", gin.H{
		// 		"title": "Hotel Management System",
		// 	})
		// })

		//Fetch all room
		authorize.GET("/api/floors", server.GetFloor)
		authorize.GET("/api/roomtypes", server.GetRoomtype)

		// //Fetch selected room information
		// defaultRoute.GET("/api/room/:roomid", server.GetService)

		/*
			---------END USER-------

		*/

		//Route for admin dashboard
		authorize.GET("/dashboard", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "dashboard.html", gin.H{
				"title": "Admin dashboard",
			})
		})

	}
}
