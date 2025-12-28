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

		/*
			---------CRUD USER-------
		*/

		//Create user route
		defaultRoute.POST("/userslist", server.CreateUser)

		//Update user route
		defaultRoute.PUT("/api/update/:userid", server.UpdateUser)

		//Delete user route
		defaultRoute.DELETE("/api/delete/:userid", server.DeleteUser)

		//Route for users list
		defaultRoute.GET("/userslist", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "users.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		//Populate user.html using api/user API
		defaultRoute.GET("/api/users", server.GetAllUsers)

		//Fetch selected user from API
		defaultRoute.GET("/api/user/:userid", server.GetUser)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD ROLE-------
		*/

		//Create role route
		defaultRoute.POST("/api/createrole", server.CreateRole)

		//Update role route
		defaultRoute.PUT("/api/updaterole/:roleid", server.UpdateRole)

		//Delete role route
		defaultRoute.DELETE("/api/deleterole/:roleid", server.DeleteRole)

		//Route for roles list
		defaultRoute.GET("/roles", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "role.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all roles
		defaultRoute.GET("/api/roles", server.GetRoles)

		//Fetch selected role information
		defaultRoute.GET("/api/roles/:roleid", server.GetRole)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD FACILITY-------
		*/
		//Create facility route
		defaultRoute.POST("/api/createfacility", server.CreateFacility)

		//Update facility route
		defaultRoute.PUT("/api/updatefacility/:facilityid", server.UpdateFacility)

		//Delete facility route
		defaultRoute.DELETE("/api/deletefacility/:facilityid", server.Deletefacility)

		//Routes for facility
		defaultRoute.GET("/facility", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "facility.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all roles
		defaultRoute.GET("/api/facility", server.GetFacilities)

		//Fetch selected role information
		defaultRoute.GET("/api/facility/:facilityid", server.GetFacility)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD SERVICE-------
		*/
		//Create services route
		defaultRoute.POST("/api/createservice", server.CreateService)

		//Update services route
		defaultRoute.PUT("/api/updateservice/:serviceid", server.UpdateService)

		//Delete services route
		defaultRoute.DELETE("/api/deleteservice/:serviceid", server.DeleteService)

		//Routes for services
		defaultRoute.GET("/service", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "service.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all services
		defaultRoute.GET("/api/services", server.GetServices)

		//Fetch selected services information
		defaultRoute.GET("/api/service/:serviceid", server.GetService)

		/*
			---------END USER-------

		*/

		/*
			---------CRUD SERVICE-------
		*/
		//Create room route
		defaultRoute.POST("/api/createroom", server.CreateRoom)

		//Update room route
		defaultRoute.PUT("/api/updateroom/:roomid", server.UpdateRoom)

		//Delete room route
		defaultRoute.DELETE("/api/deleteroom/:roomid", server.DeleteRoom)

		//Routes for room
		defaultRoute.GET("/rooms", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "room.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Fetch all room
		defaultRoute.GET("/api/rooms", server.GetRooms)

		//Fetch selected room information
		defaultRoute.GET("/api/room/:roomid", server.GetRoom)

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
		defaultRoute.GET("/api/floors", server.GetFloor)
		defaultRoute.GET("/api/roomtypes", server.GetRoomtype)

		// //Fetch selected room information
		// defaultRoute.GET("/api/room/:roomid", server.GetService)

		/*
			---------END USER-------

		*/
	}
}
