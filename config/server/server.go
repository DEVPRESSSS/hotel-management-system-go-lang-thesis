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

	//router.LoadHTMLGlob("views/**/*.html")
	// Register routes

	files := []string{
		//Header and footer for default page
		"views/Layout/footer.html",
		"views/Layout/header.html",
		//Error html page
		"views/ErrorView/errors.html",
		//Default page
		"views/defaultView/index.html",
		//User html
		"views/Areas/Admin/users/users.html",
		//Roles html
		"views/Areas/Admin/roles/role.html",
		//Facility html
		"views/Areas/Admin/facilities/facility.html",
		//Service html
		"views/Areas/Admin/services/service.html",
		//Room html
		"views/Areas/Admin/rooms/room.html",
		//Login html
		"views/Auth/login.html",
		//Admin dashboard layout
		"views/Areas/Admin/dashboard_header.html",
		"views/Areas/Admin/dashboard/dashboard.html",
		"views/Areas/Admin/dashboard_footer.html",
	}
	router.LoadHTMLFiles(files...)

	routes.AuthRoutes(db, router)
	router.Run(host)
}
