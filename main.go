package main

import (
	"log"
	"os"

	"HMS-GO/config/database"
	"HMS-GO/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := models.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	_ = db

	r := gin.Default()
	r.Run(":8085")
}
