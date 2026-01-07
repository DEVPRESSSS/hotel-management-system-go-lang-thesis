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

		// Public
		"views/defaultView/index.html",
		"views/Auth/login.html",

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
	}
	router.LoadHTMLFiles(files...)

	routes.AuthRoutes(db, router)
	router.Run(host)
}
