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
		"views/Layout/footer.html",
		"views/Layout/header.html",
		"views/ErrorView/errors.html",
		"views/defaultView/index.html",
		"views/Areas/Admin/users/users.html",
		"views/Areas/Admin/roles/role.html",
		"views/Areas/Admin/facilities/facility.html",
	}
	//Try to load all the files here using LoadHtmlfiles not HtmlGlob
	router.LoadHTMLFiles(files...)

	routes.AuthRoutes(db, router)
	router.Run(host)
}
