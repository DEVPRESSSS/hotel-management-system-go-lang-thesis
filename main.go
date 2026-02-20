package main

import (
	"log"
	"os"

	"HMS-GO/config/database"
	"HMS-GO/config/server"
)

func main() {

	// godotenv.Load()
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	host := os.Getenv("HOST_ADDR")
	if host == "" {
		host = ":8085"
	}

	// cfg := models.DatabaseConfig{
	// 	Host:     os.Getenv("dbHost"),
	// 	Port:     os.Getenv("dbPort"),
	// 	User:     os.Getenv("dbUser"),
	// 	Password: os.Getenv("dbPassword"),
	// 	DBName:   os.Getenv("dbName"),
	// }

	db, err := database.InitDatabase()
	database.SeedRoles(db)
	database.SeedAccess(db)
	database.SeedRoleAccess(db)
	database.SeedAdminUser(db)
	database.SeedRoomTypes(db)
	database.SeedFloors(db)
	database.SeedFoodCategories(db)

	if err != nil {

		log.Fatal("Database connection failed:", err)
	}

	_ = db

	server.SetupServer(host, db)

}
