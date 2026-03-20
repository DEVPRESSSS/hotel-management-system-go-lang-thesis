package server

import (
	"HMS-GO/internal/routes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func loadTemplates(router *gin.Engine) {
	tmpl := template.New("")

	err := filepath.WalkDir("views", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".html" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Register template using its full path as the name
			template.Must(tmpl.New(path).Parse(string(content)))
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk template directory: %v", err)
	}

	router.SetHTMLTemplate(tmpl)
}

func SetupServer(host string, db *gorm.DB) {
	router := gin.Default()

	router.SetTrustedProxies([]string{"localhost"})

	router.Static("/src", "./src")
	router.Static("/food_images", "./src/food_images")
	router.Static("/room_images", "./src/room_images")

	loadTemplates(router)

	routes.AuthRoutes(db, router)
	router.Run(host)
}

// func SetupServer(host string, db *gorm.DB) {
// 	router := gin.Default()

// 	router.SetTrustedProxies([]string{"localhost"})

// 	router.Static("/src", "./src")
// 	router.Static("/food_images", "./src/food_images")
// 	router.Static("/room_images", "./src/room_images")
// 	files := []string{
// 		// Default layout
// 		"views/Layout/header.html",
// 		"views/Layout/footer.html",

// 		// Error
// 		"views/ErrorView/errors.html",
// 		"views/ErrorView/404.html",
// 		"views/ErrorView/forbidden.html",

// 		// Public
// 		"views/defaultView/index.html",

// 		"views/Auth/login.html",
// 		"views/Auth/register.html",
// 		"views/Auth/forgot_password.html",
// 		"views/Auth/resetPassword.html",
// 		"views/Auth/reset_password_form.html",

// 		//Guest layout
// 		"views/Areas/Guest/guest_header.html",
// 		"views/Areas/Guest/guest_footer.html",
// 		"views/Areas/Guest/dashboard/guest_index.html",
// 		"views/Areas/Guest/booking/booking.html",
// 		"views/Areas/Guest/booking/confirm_booking.html",
// 		"views/Areas/Guest/booking/booking_success.html",
// 		"views/Areas/Guest/request/food_payment_success.html",
// 		"views/Areas/Guest/booking/booking_history.html",
// 		"views/Areas/Guest/request/food_request.html",

// 		// Admin layouts
// 		"views/Areas/Admin/Layout.html",
// 		"views/Areas/Admin/dashboard_header.html",
// 		"views/Areas/Admin/dashboard_footer.html",

// 		// Admin pages
// 		"views/Areas/Admin/dashboard/dashboard.html",
// 		"views/Areas/Admin/users/users.html",
// 		"views/Areas/Admin/roles/role.html",
// 		"views/Areas/Admin/facilities/facility.html",
// 		"views/Areas/Admin/services/service.html",
// 		"views/Areas/Admin/services/food_service.html",
// 		"views/Areas/Admin/rooms/room.html",
// 		"views/Areas/Admin/rbac/rbac.html",
// 		"views/Areas/Admin/aminity/aminity.html",
// 		"views/Areas/Admin/roomaminity/room_aminity.html",
// 		"views/Areas/Admin/roomtype/roomtype.html",
// 		"views/Areas/Admin/guests/guest.html",
// 		"views/Areas/Admin/reservation/reservation.html",
// 		"views/Areas/Admin/reservation/walkin_booking.html",
// 		"views/Areas/Admin/reservation/walkin_room_details.html",
// 		"views/Areas/Admin/reservation/walkin_confirm_booking.html",
// 		"views/Areas/Admin/calendar/calendar.html",
// 		"views/Areas/Admin/settings/logs.html",
// 		"views/Areas/Admin/maintenance/maintenance.html",
// 		"views/Areas/Admin/2d/2d.html",
// 		"views/Areas/Admin/services/food_category.html",
// 	}
// 	router.LoadHTMLFiles(files...)

// 	routes.AuthRoutes(db, router)
// 	router.Run(host)
// }
