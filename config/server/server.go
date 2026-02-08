package server

import (
	"HMS-GO/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupServer(host string, db *gorm.DB) {
	router := gin.Default()

	router.SetTrustedProxies([]string{"localhost"})

	router.Static("/src", "./src")

	files := []string{
		// Default layout
		"views/Layout/header.html",
		"views/Layout/footer.html",

		// Error
		"views/ErrorView/errors.html",
		"views/ErrorView/404.html",
		"views/ErrorView/forbidden.html",

		// Public
		"views/defaultView/index.html",

		"views/Auth/login.html",
		"views/Auth/register.html",
		"views/Auth/forgot_password.html",
		"views/Auth/resetPassword.html",
		"views/Auth/reset_password_form.html",

		//Guest layout
		"views/Areas/Guest/guest_header.html",
		"views/Areas/Guest/guest_footer.html",
		"views/Areas/Guest/dashboard/guest_index.html",
		"views/Areas/Guest/booking/booking.html",
		"views/Areas/Guest/booking/confirm_booking.html",
		"views/Areas/Guest/booking/booking_success.html",

		// Admin layouts
		"views/Areas/Admin/Layout.html",
		"views/Areas/Admin/dashboard_header.html",
		"views/Areas/Admin/dashboard_footer.html",

		// Admin pages
		"views/Areas/Admin/dashboard/dashboard.html",
		"views/Areas/Admin/users/users.html",
		"views/Areas/Admin/roles/role.html",
		"views/Areas/Admin/facilities/facility.html",
		"views/Areas/Admin/services/service.html",
		"views/Areas/Admin/rooms/room.html",
		"views/Areas/Admin/rbac/rbac.html",
		"views/Areas/Admin/aminity/aminity.html",
		"views/Areas/Admin/roomaminity/room_aminity.html",
		"views/Areas/Admin/roomtype/roomtype.html",
		"views/Areas/Admin/guests/guest.html",
		"views/Areas/Admin/reservation/reservation.html",
		"views/Areas/Admin/reservation/walkin_booking.html",
		"views/Areas/Admin/reservation/walkin_room_details.html",
		"views/Areas/Admin/reservation/walkin_confirm_booking.html",
		"views/Areas/Admin/calendar/calendar.html",
		"views/Areas/Admin/settings/logs.html",
	}
	router.LoadHTMLFiles(files...)

	routes.AuthRoutes(db, router)
	router.Run(host)
}
