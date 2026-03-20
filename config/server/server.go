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
