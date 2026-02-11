package routes

import (
	"HMS-GO/internal/controllers"
	"HMS-GO/internal/middlewares"
	rbac "HMS-GO/internal/middlewares/RBAC"
	"HMS-GO/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(db *gorm.DB, router *gin.Engine) {

	// ==================================================================================================================
	// GLOBAL HANDLERS
	// ==================================================================================================================

	// Custom 404 handler
	router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{})
	})

	// Initialize controller server
	server := controllers.NewServer(db)

	// ==================================================================================================================
	// PUBLIC ROUTES (NO AUTHENTICATION REQUIRED)
	// ==================================================================================================================
	defaultRoute := router.Group("/")
	{
		// Home page
		defaultRoute.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Login routes
		defaultRoute.POST("/api/auth", server.Login)
		defaultRoute.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "login.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Registration routes
		defaultRoute.POST("/createaccount", server.RegisterGuest)
		defaultRoute.GET("/register", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "register.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Password recovery routes
		defaultRoute.POST("/verifyemail", server.RegisterGuest)
		defaultRoute.GET("/forgotpassword", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "forgot_password.html", gin.H{
				"title": "Hotel Management System",
			})
		})
		defaultRoute.POST("/forgotpassword", server.ForgotPassword)
		defaultRoute.PATCH("/api/resetpassword/:resetToken", server.ResetPassword)
		defaultRoute.GET("/resetpassword/:resetToken", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "resetPassword.html", nil)
		})
		defaultRoute.GET("/resetpassword-form/:resetToken", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "reset_password_form.html", nil)
		})
	}

	// ==================================================================================================================
	// PROTECTED ROUTES (JWT AUTHENTICATION REQUIRED)
	// ==================================================================================================================
	authorize := router.Group("/")
	authorize.Use(middlewares.AuthMiddleware())
	{

		// ==============================================================================================================
		// ADMIN DASHBOARD
		// ==============================================================================================================
		authorize.GET("/api/dashboard", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "dashboard.html", gin.H{
				"title": "Dashboard",
			})
		})

		// ==============================================================================================================
		// GUEST ROUTES
		// ==============================================================================================================

		// Guest dashboard
		authorize.GET("/guest/dashboard", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "guest_index.html", gin.H{
				"title": "Hotel Management System",
			})
		})

		// Get available rooms
		authorize.GET("/avail/rooms", rbac.RBACMiddleware("booking"), server.GetRooms)

		// Calendar and room selection
		authorize.GET("/api/calendar/:room_id", rbac.RBACMiddleware("booking"), server.FetchCalendar)
		authorize.GET("/api/roomselected/:roomid", rbac.RBACMiddleware("booking"), server.RoomSelected)

		// Room details page
		authorize.GET("/roomdetails", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "booking.html", gin.H{
				"title": "Room details",
			})
		})

		// Booking and payment
		authorize.POST("/api/booking/calculate", rbac.RBACMiddleware("booking"), server.CalculateBookingPrice)
		authorize.POST("/api/create-checkout-session", rbac.RBACMiddleware("booking"), server.CreateCheckoutSession)
		authorize.POST("/api/booking/confirmbooking", rbac.RBACMiddleware("booking"), server.ConfirmBooking)
		authorize.GET("/booking/success", server.BookingSuccess)
		authorize.POST("/paymongo/create/payment-intent", server.CreatePaymentIntent)

		// Booking summary page
		authorize.GET("/booking/summary", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "confirm_booking.html", gin.H{
				"title": "Room details",
			})
		})

		// ==============================================================================================================
		// GUEST MANAGEMENT
		// ==============================================================================================================

		// Render guest page
		authorize.GET("/guest", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "guest.html", gin.H{
				"title": "Guest",
			})
		})

		// Guest APIs
		authorize.GET("/api/guest", rbac.RBACMiddleware("read"), server.GetAllGuest)

		// ==============================================================================================================
		// RESERVATION MANAGEMENT
		// ==============================================================================================================

		// Reservation APIs
		authorize.GET("/api/reservations", rbac.RBACMiddleware("read"), server.GetAllReservations)
		authorize.GET("/api/reservations/clean/:id", rbac.RBACMiddleware("read"), server.GetAllReservations)
		authorize.POST("/api/reservations/checkin/:id", rbac.RBACMiddleware("create"), server.CheckinStatus)
		authorize.GET("/api/reservations/events", rbac.RBACMiddleware("read"), server.GetAllEventsReservations)
		authorize.POST("/api/reservations/assigncleaner/:id", rbac.RBACMiddleware("read"), server.AssignCleaner)

		// Render reservation page
		authorize.GET("/reservations", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "reservation.html", gin.H{
				"title": "Reservation",
			})
		})

		// Calendar page
		authorize.GET("/reservation-calendar", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "calendar.html", gin.H{
				"title": "Calendar",
			})
		})

		// Walk-in booking pages
		authorize.GET("/api/walkin-booking", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "walkin_booking.html", gin.H{
				"title": "Booking",
			})
		})

		authorize.GET("/api/walkin/room-details", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "walkin_room_details.html", gin.H{
				"title": "Booking",
			})
		})

		authorize.GET("/api/walkin/confirm-booking", rbac.RBACMiddleware("booking"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "walkin_confirm_booking.html", gin.H{
				"title": "Booking",
			})
		})

		// ==============================================================================================================
		// MAINTENANCE MANAGEMENT
		// ==============================================================================================================

		// Maintenance APIs
		authorize.GET("/api/maintenances", rbac.RBACMiddleware("read"), server.GetAllCleaners)
		authorize.GET("/api/maintenances/:id", rbac.RBACMiddleware("read"), server.GetCleaner)
		authorize.POST("/api/maintenances", rbac.RBACMiddleware("create"), server.CreateCleaner)
		authorize.PUT("/api/maintenances/:id", rbac.RBACMiddleware("update"), server.UpdateCleaner)
		authorize.DELETE("/api/maintenances/:id", rbac.RBACMiddleware("delete"), server.DeleteCleaner)

		// Render maintenance page
		authorize.GET("/maintenance", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "maintenance.html", gin.H{
				"title": "Maintenance",
			})
		})

		// ==============================================================================================================
		// ROOM MANAGEMENT
		// ==============================================================================================================

		// Room APIs
		authorize.GET("/api/rooms", rbac.RBACMiddleware("read"), server.GetRooms)
		authorize.GET("/api/room/:roomid", rbac.RBACMiddleware("read"), server.GetRoom)
		authorize.POST("/api/createroom", rbac.RBACMiddleware("create"), server.CreateRoom)
		authorize.PUT("/api/updateroom/:roomid", rbac.RBACMiddleware("update"), server.UpdateRoom)
		authorize.DELETE("/api/deleteroom/:roomid", rbac.RBACMiddleware("delete"), server.DeleteRoom)

		// Render rooms page
		authorize.GET("/rooms", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "room.html", gin.H{
				"title": "Room",
			})
		})

		// ==============================================================================================================
		// ROOM TYPE MANAGEMENT
		// ==============================================================================================================

		// Room Type APIs
		authorize.GET("/api/roomtypes", rbac.RBACMiddleware("read"), server.GetRoomtype)
		authorize.GET("/api/roomtype/:roomtypeid", rbac.RBACMiddleware("read"), server.GetRoomTypeRecord)
		authorize.POST("/api/createroomtype", rbac.RBACMiddleware("create"), server.CreateRoomType)
		authorize.PUT("/api/updateroomtype/:roomtypeid", rbac.RBACMiddleware("update"), server.UpdateRoomType)
		authorize.DELETE("/api/deleteroomtype/:roomtypeid", rbac.RBACMiddleware("delete"), server.DeleteRoomType)

		// Render room type page
		authorize.GET("/roomtype", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "roomtype.html", gin.H{
				"title": "Room Type",
			})
		})

		// ==============================================================================================================
		// AMENITY MANAGEMENT
		// ==============================================================================================================

		// Amenity APIs
		authorize.GET("/api/aminities", rbac.RBACMiddleware("read"), server.GetAminities)
		authorize.GET("/api/aminity/:amenityid", rbac.RBACMiddleware("read"), server.GetAminity)
		authorize.POST("/api/createaminity", rbac.RBACMiddleware("create"), server.CreateAminity)
		authorize.PUT("/api/updateaminity/:amenityid", rbac.RBACMiddleware("update"), server.UpdateAminity)
		authorize.DELETE("/api/deleteamenity/:amenityid", rbac.RBACMiddleware("delete"), server.DeleteAminity)

		// Render amenities page
		authorize.GET("/amenities", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "aminity.html", gin.H{
				"title": "Amenities",
			})
		})

		// ==============================================================================================================
		// ROOM AMENITY MANAGEMENT
		// ==============================================================================================================

		// Room Amenity APIs
		authorize.GET("/api/roomaminities", rbac.RBACMiddleware("read"), server.GetRoomAminities)
		authorize.GET("/api/roomaminity/:roomid", rbac.RBACMiddleware("read"), server.GetAminity)
		authorize.POST("/api/createroomaminity", rbac.RBACMiddleware("create"), server.CreateRoomAminity)
		authorize.PUT("/api/updateroomaminity/:roomid", rbac.RBACMiddleware("update"), server.UpdateRoomAminity)
		authorize.DELETE("/api/deleteroomaminity/:roomid", rbac.RBACMiddleware("delete"), server.DeleteRoomAminity)

		// Render room amenities page
		authorize.GET("/roomaminities", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "room_aminity.html", gin.H{
				"title": "Room Amenities",
			})
		})

		// ==============================================================================================================
		// SERVICE MANAGEMENT
		// ==============================================================================================================

		// Service APIs
		authorize.GET("/api/services", rbac.RBACMiddleware("read"), server.GetServices)
		authorize.GET("/api/service/:serviceid", rbac.RBACMiddleware("read"), server.GetService)
		authorize.POST("/api/createservice", rbac.RBACMiddleware("create"), server.CreateService)
		authorize.PUT("/api/updateservice/:serviceid", rbac.RBACMiddleware("update"), server.UpdateService)
		authorize.DELETE("/api/deleteservice/:serviceid", rbac.RBACMiddleware("delete"), server.DeleteService)

		// Render services page
		authorize.GET("/service", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "service.html", gin.H{
				"title": "Services",
			})
		})

		// ==============================================================================================================
		// FACILITY MANAGEMENT
		// ==============================================================================================================

		// Facility APIs
		authorize.GET("/api/facility", rbac.RBACMiddleware("read"), server.GetFacilities)
		authorize.GET("/api/facility/:facilityid", rbac.RBACMiddleware("read"), server.GetFacility)
		authorize.POST("/api/createfacility", rbac.RBACMiddleware("create"), server.CreateFacility)
		authorize.PUT("/api/updatefacility/:facilityid", rbac.RBACMiddleware("update"), server.UpdateFacility)
		authorize.DELETE("/api/deletefacility/:facilityid", rbac.RBACMiddleware("delete"), server.Deletefacility)

		// Render facility page
		authorize.GET("/facility", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "facility.html", gin.H{
				"title": "Facility",
			})
		})

		// ==============================================================================================================
		// FLOOR MANAGEMENT
		// ==============================================================================================================

		// Floor APIs
		authorize.GET("/api/floors", rbac.RBACMiddleware("read"), server.GetFloor)

		// ==============================================================================================================
		// USER MANAGEMENT
		// ==============================================================================================================

		// User APIs
		authorize.GET("/api/users", rbac.RBACMiddleware("read"), server.GetAllUsers)
		authorize.GET("/api/user/:userid", rbac.RBACMiddleware("read"), server.GetUser)
		authorize.POST("/userslist", rbac.RBACMiddleware("create"), server.CreateUser)
		authorize.PUT("/api/update/:userid", rbac.RBACMiddleware("update"), server.UpdateUser)
		authorize.DELETE("/api/delete/:userid", rbac.RBACMiddleware("delete"), server.DeleteUser)

		// Render users page
		authorize.GET("/users", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "users.html", gin.H{
				"title": "Users",
			})
		})

		// ==============================================================================================================
		// ROLE MANAGEMENT
		// ==============================================================================================================

		// Role APIs
		authorize.GET("/api/roles", rbac.RBACMiddleware("read"), server.GetRoles)
		authorize.GET("/api/roles/:roleid", rbac.RBACMiddleware("read"), server.GetRole)
		authorize.POST("/api/createrole", rbac.RBACMiddleware("create"), server.CreateRole)
		authorize.PUT("/api/updaterole/:roleid", rbac.RBACMiddleware("update"), server.UpdateRole)
		authorize.DELETE("/api/deleterole/:roleid", rbac.RBACMiddleware("delete"), server.DeleteRole)

		// Render roles page
		authorize.GET("/roles", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "role.html", gin.H{
				"title": "Roles",
			})
		})

		// ==============================================================================================================
		// ROLE BASED ACCESS CONTROL (RBAC)
		// ==============================================================================================================

		// RBAC APIs
		authorize.GET("/api/rbac", rbac.RBACMiddleware("read"), server.RoleAccess)
		authorize.POST("/api/createrc", server.CreateRoleAccess)
		authorize.POST("/api/updaterc/:roleid", server.UpdateRoleAcccess)
		authorize.DELETE("/api/deleterc/:accessid", server.DeleteRoleAccess)

		// Render RBAC page
		authorize.GET("/rbac", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "rbac.html", gin.H{
				"title": "Role Based Access",
			})
		})

		// ==============================================================================================================
		// ACCESS MANAGEMENT
		// ==============================================================================================================

		// Access APIs
		authorize.GET("/api/access", rbac.RBACMiddleware("read"), server.Access)
		authorize.GET("/api/access/:accessid", rbac.RBACMiddleware("read"), server.GetAccess)
		authorize.POST("/api/createac", rbac.RBACMiddleware("create"), server.CreateAccess)
		authorize.POST("/api/updateac/:accessid", rbac.RBACMiddleware("update"), server.UpdateAcccess)
		authorize.DELETE("/api/deleteac/:accessid", rbac.RBACMiddleware("delete"), server.DeleteAccess)

		// ==============================================================================================================
		// ACTIVITY LOGS
		// ==============================================================================================================

		// Activity Logs APIs
		authorize.GET("/api/getlogs/", rbac.RBACMiddleware("read"), server.GetLogs)

		// Render logs page
		authorize.GET("/logs", rbac.RBACMiddleware("read"), func(ctx *gin.Context) {
			utils.RenderWithRole(ctx, "logs.html", gin.H{
				"title": "Activity Logs",
			})
		})

		// ==============================================================================================================
		// AUTHENTICATION
		// ==============================================================================================================

		// Logout routes
		authorize.POST("/logout", server.Logout)
		authorize.GET("/logout", server.Logout)
	}
}
