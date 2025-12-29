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

	type User struct {
		UserId   string
		Username string
		Email    string
	}

	// token, err := utils.CreateToken(models.User{
	// 	UserId:   "USER-101",
	// 	Username: "Jerald",
	// 	Email:    "xmont@gmail.com",
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	//fmt.Println(token)

	server.SetupServer(host, db)

}
