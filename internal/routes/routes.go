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

		// Login page
		defaultRoute.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		defaultRoute.POST("/createaccount", server.RegisterGuest)
		//Register page
		defaultRoute.GET("/register", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "register.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		defaultRoute.POST("/verifyemail", server.RegisterGuest)
		//Register page
		defaultRoute.GET("/forgotpassword", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "forgot_password.html", gin.H{
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

		//GUEST ROUTES
		//===============================================
		// Guest dashboard
		authorize.GET("/guest/dashboard", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "guest_index.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Get rooms
		authorize.GET("/avail/rooms", rbac.RBACMiddleware("booking"), server.GetRooms)
		//Api for getting all the reservations in calendar
		authorize.GET("/api/calendar/:room_id", rbac.RBACMiddleware("booking"), server.FetchCalendar)
		//Get the selected room
		authorize.GET("/api/roomselected/:roomid", rbac.RBACMiddleware("booking"), server.RoomSelected)
		//Populate the details of the room
		authorize.GET("/roomdetails", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "booking.html", gin.H{
				"title": "Room details",
			})
		})
		//Calculate the booking price per guest
		authorize.POST("/api/booking/calculate", rbac.RBACMiddleware("booking"), server.CalculateBookingPrice)
		authorize.POST("/api/create-checkout-session", rbac.RBACMiddleware("booking"), server.CreateCheckoutSession)
		authorize.POST("/api/booking/confirmbooking", rbac.RBACMiddleware("booking"), server.ConfirmBooking)
		authorize.GET("/booking/success", server.BookingSuccess)
		authorize.POST("/paymongo/create/payment-intent", server.CreatePaymentIntent)
		//Render confirm booking html
		authorize.GET("/booking/summary", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "confirm_booking.html", gin.H{
				"title": "Room details",
			})
		})

		// ==============================================
		// RESERVATION  MANAGEMENT
		// ==============================================

		//Api for getting all the reservations
		authorize.GET("/api/reservations", rbac.RBACMiddleware("read"), server.GetAllReservations)

		//Render the reservation page
		authorize.GET("/reservations", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "reservation.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		//Walkin booking
		authorize.GET("/api/walkin-booking", rbac.RBACMiddleware("create"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "walkin_booking.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		//Walkin booking room details
		authorize.GET("/api/walkin/room-details", rbac.RBACMiddleware("create"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "walkin_room_details.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		//Walkin booking room details
		authorize.GET("/api/walkin/confirm-booking", rbac.RBACMiddleware("create"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "walkin_confirm_booking.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		authorize.GET("/api/reservations/events", rbac.RBACMiddleware("read"), server.GetAllEventsReservations)
		authorize.GET("/reservation-calendar", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "calendar.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// ==============================================
		// GUEST  MANAGEMENT
		// ==============================================

		// Render guest page
		authorize.GET("guest", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "guest.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// ==============================================
		// ROOM TYPES MANAGEMENT (CRUD)
		// ==============================================

		//Create room type route
		authorize.POST("/api/createroomtype", rbac.RBACMiddleware("create"), server.CreateRoomType)
		//Update room type route
		authorize.PUT("/api/updateroomtype/:roomtypeid", rbac.RBACMiddleware("update"), server.UpdateRoomType)
		//Delete room type route
		authorize.DELETE("/api/deleteroomtype/:roomtypeid", rbac.RBACMiddleware("delete"), server.DeleteRoomType)
		// Room Types page
		authorize.GET("/roomtype", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "roomtype.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		//Room type APIs
		authorize.GET("/api/roomtypes", rbac.RBACMiddleware("read"), server.GetRoomtype)
		//fetch one room type
		authorize.GET("/api/roomtype/:roomtypeid", rbac.RBACMiddleware("read"), server.GetRoomTypeRecord)

		// ==============================================
		// ROOM AMENITY MANAGEMENT (CRUD)
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
		// AMENITY MANAGEMENT (CRUD)
		// ==============================================

		authorize.POST("/api/createaminity", rbac.RBACMiddleware("create"), server.CreateAminity)
		authorize.PUT("/api/updateaminity/:amenityid", rbac.RBACMiddleware("update"), server.UpdateAminity)
		authorize.DELETE("/api/deleteamenity/:amenityid", rbac.RBACMiddleware("delete"), server.DeleteAminity)

		// Amenity
		authorize.GET("/amenities", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "aminity.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Aminity APIs
		authorize.GET("/api/aminities", rbac.RBACMiddleware("read"), server.GetAminities)
		authorize.GET("/api/aminity/:amenityid", rbac.RBACMiddleware("read"), server.GetAminity)

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
		authorize.GET("/api/guest", rbac.RBACMiddleware("read"), server.GetAllGuest)
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

		authorize.GET("/api/dashboard", rbac.RBACMiddleware("create"), rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
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

		//Logout
		authorize.POST("/logout", server.Logout)
		authorize.GET("/logout", server.Logout)

	}
}
