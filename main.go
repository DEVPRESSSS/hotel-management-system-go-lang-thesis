package main

import (
	"log"
	"os"

	"HMS-GO/config/database"
	"HMS-GO/config/server"
	"HMS-GO/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	//router := gin.Default()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST_ADDR")
	if host == "" {
		host = ":8085"
	}

	cfg := models.DatabaseConfig{
		Host:     os.Getenv("dbHost"),
		Port:     os.Getenv("dbPort"),
		User:     os.Getenv("dbUser"),
		Password: os.Getenv("dbPassword"),
		DBName:   os.Getenv("dbName"),
	}

	db, err := database.InitDatabase(cfg)

	if err != nil {

		log.Fatal("Database connection failed:", err)
	}

	_ = db

	server.SetupServer(host, db)

}
