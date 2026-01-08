package routes

import (
	"HMS-GO/internal/controllers"
	"HMS-GO/internal/middlewares"
	rbac "HMS-GO/internal/middlewares/RBAC"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(db *gorm.DB, router *gin.Engine) {

	// --------------------------------------------------
	// Global Handlers
	// --------------------------------------------------

	// Custom 404 handler
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{})
	})

	// Initialize controller server
	server := controllers.NewServer(db)

	// --------------------------------------------------
	// Public Routes (No Authentication Required)
	// --------------------------------------------------
	defaultRoute := router.Group("/")
	{
		// Home page
		defaultRoute.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		//Login post method handler
		defaultRoute.POST("/api/auth", server.Login)
		// Get rooms
		defaultRoute.GET("/avail/rooms", server.GetRooms)
		//Get the selected room
		defaultRoute.GET("/api/roomselected/:roomid", server.RoomSelected)

		// Login page
		defaultRoute.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"title": "Hotel Management System",
			})
		})
	}

	// --------------------------------------------------
	// Protected Routes (JWT Authentication Required)
	// --------------------------------------------------
	authorize := router.Group("/")
	authorize.Use(middlewares.AuthMiddleware())
	{
		// ==============================================
		// ROOM TYPES MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/api/createroomtype", rbac.RBACMiddleware("create"), server.CreateRoomType)
		authorize.PUT("/api/updateroomtype/:roomtypeid", rbac.RBACMiddleware("update"), server.UpdateRoomType)
		authorize.DELETE("/api/deleteroomtype/:roomtypeid", rbac.RBACMiddleware("delete"), server.DeleteRoomType)

		// Room Types page
		authorize.GET("/roomtype", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "roomtype.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Room type APIs
		authorize.GET("/api/roomtypes", rbac.RBACMiddleware("read"), server.GetRoomtype)
		authorize.GET("/api/roomtype/:roomtypeid", rbac.RBACMiddleware("read"), server.GetRoomTypeRecord)

		// ==============================================
		// ROOM AMINITY MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/api/createroomaminity", rbac.RBACMiddleware("create"), server.CreateRoomAminity)
		authorize.PUT("/api/updateroomaminity/:roomid", rbac.RBACMiddleware("update"), server.UpdateRoomAminity)
		authorize.DELETE("/api/deleteroomaminity/:roomid", rbac.RBACMiddleware("delete"), server.DeleteRoomAminity)

		// Room Aminity
		authorize.GET("/roomaminities", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "room_aminity.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		//Room Aminity APIs
		authorize.GET("/api/roomaminities", rbac.RBACMiddleware("read"), server.GetRoomAminities)
		authorize.GET("/api/roomaminity/:roomid", rbac.RBACMiddleware("read"), server.GetAminity)

		// ==============================================
		// AMINITY MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/api/createaminity", rbac.RBACMiddleware("create"), server.CreateAminity)
		authorize.PUT("/api/updateaminity/:aminityid", rbac.RBACMiddleware("update"), server.UpdateAminity)
		authorize.DELETE("/api/deleteaminity/:aminityid", rbac.RBACMiddleware("delete"), server.DeleteAminity)

		// Aminity
		authorize.GET("/aminities", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "aminity.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Aminity APIs
		authorize.GET("/api/aminities", rbac.RBACMiddleware("read"), server.GetAminities)
		authorize.GET("/api/aminity/:aminityid", rbac.RBACMiddleware("read"), server.GetAminity)

		// ==============================================
		// USER MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/userslist", rbac.RBACMiddleware("create"), server.CreateUser)
		authorize.PUT("/api/update/:userid", rbac.RBACMiddleware("update"), server.UpdateUser)
		authorize.DELETE("/api/delete/:userid", rbac.RBACMiddleware("delete"), server.DeleteUser)

		// Users page
		authorize.GET("/users", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "users.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// User APIs
		authorize.GET("/api/users", rbac.RBACMiddleware("read"), server.GetAllUsers)
		authorize.GET("/api/user/:userid", rbac.RBACMiddleware("read"), server.GetUser)

		// ==============================================
		// ROLE MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/api/createrole", rbac.RBACMiddleware("create"), server.CreateRole)
		authorize.PUT("/api/updaterole/:roleid", rbac.RBACMiddleware("update"), server.UpdateRole)
		authorize.DELETE("/api/deleterole/:roleid", rbac.RBACMiddleware("delete"), server.DeleteRole)

		// Roles page
		authorize.GET("/roles", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "role.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Role APIs
		authorize.GET("/api/roles", rbac.RBACMiddleware("read"), server.GetRoles)
		authorize.GET("/api/roles/:roleid", rbac.RBACMiddleware("read"), server.GetRole)

		// ==============================================
		// FACILITY MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/api/createfacility", rbac.RBACMiddleware("create"), server.CreateFacility)
		authorize.PUT("/api/updatefacility/:facilityid", rbac.RBACMiddleware("update"), server.UpdateFacility)
		authorize.DELETE("/api/deletefacility/:facilityid", rbac.RBACMiddleware("delete"), server.Deletefacility)

		// Facility page (RBAC: Read Permission)
		authorize.GET("/facility", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "facility.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Facility APIs
		authorize.GET("/api/facility", rbac.RBACMiddleware("read"), server.GetFacilities)
		authorize.GET("/api/facility/:facilityid", rbac.RBACMiddleware("read"), server.GetFacility)

		// ==============================================
		// SERVICE MANAGEMENT (CRUD + RBAC)
		// ==============================================

		authorize.POST("/api/createservice", rbac.RBACMiddleware("create"), server.CreateService)
		authorize.PUT("/api/updateservice/:serviceid", rbac.RBACMiddleware("update"), server.UpdateService)
		authorize.DELETE("/api/deleteservice/:serviceid", rbac.RBACMiddleware("delete"), server.DeleteService)

		// Services page
		authorize.GET("/service", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "service.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Service APIs
		authorize.GET("/api/services", rbac.RBACMiddleware("read"), server.GetServices)
		authorize.GET("/api/service/:serviceid", rbac.RBACMiddleware("read"), server.GetService)

		// ==============================================
		// ROOM MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/api/createroom", rbac.RBACMiddleware("create"), server.CreateRoom)
		authorize.PUT("/api/updateroom/:roomid", rbac.RBACMiddleware("update"), server.UpdateRoom)
		authorize.DELETE("/api/deleteroom/:roomid", rbac.RBACMiddleware("delete"), server.DeleteRoom)

		// Rooms page
		authorize.GET("/rooms", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "room.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Room APIs
		authorize.GET("/api/rooms", rbac.RBACMiddleware("read"), server.GetRooms)
		authorize.GET("/api/room/:roomid", rbac.RBACMiddleware("read"), server.GetRoom)

		// ==============================================
		// FLOOR & ROOM TYPE (READ ONLY)
		// ==============================================

		authorize.GET("/api/floors", rbac.RBACMiddleware("read"), server.GetFloor)

		// ==============================================
		// ADMIN DASHBOARD (RBAC Protected)
		// ==============================================

		authorize.GET("/api/dashboard", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "dashboard.html", gin.H{
				"title": "Admin Dashboard",
			})
		})

		// ==============================================
		// ROLE BASED ACCESS ROUTE
		// ==============================================

		//Create role access
		authorize.POST("/api/createrc", server.CreateRoleAccess)
		//Update role access
		authorize.POST("/api/updaterc/:roleid", server.UpdateRoleAcccess)
		//Delete role access
		authorize.DELETE("/api/deleterc/:accessid", server.DeleteRoleAccess)
		//Role based access api
		authorize.GET("/api/rbac", rbac.RBACMiddleware("read"), server.RoleAccess)
		// ==============================================
		// ACCESS ROUTE
		// ==============================================

		//Create access function
		authorize.POST("/api/createac", rbac.RBACMiddleware("create"), server.CreateAccess)
		//Update access function
		authorize.POST("/api/updateac/:accessid", rbac.RBACMiddleware("update"), server.UpdateAcccess)
		//Delete access function
		authorize.DELETE("/api/deleteac/:accessid", rbac.RBACMiddleware("delete"), server.DeleteAccess)
		//Get access record
		authorize.GET("/api/access/:accessid", rbac.RBACMiddleware("read"), server.GetAccess)
		//Get all access
		authorize.GET("/api/access", rbac.RBACMiddleware("read"), server.Access)
		authorize.GET("/rbac", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "rbac.html", gin.H{
				"title": "Hotel Management System",
			})
		})

	}
}
