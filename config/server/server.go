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

	router.LoadHTMLGlob("views/**/*.html")
	// Register routes
	routes.AuthRoutes(db, router)

	router.Run(host)
}
