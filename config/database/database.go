package database

import (
	"HMS-GO/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(cfg models.DatabaseConfig) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Access{},
		&models.RoleAccess{},
		&models.UserRole{},
		&models.Facility{},
		&models.Service{},
		&models.Floor{},
		&models.RoomType{},
		&models.Room{},
		&models.Amenity{},
		&models.RoomAmenity{},
		&models.Book{},
		&models.BookingGuest{},
		&models.HistoryLog{},
		&models.Cleaner{},
		&models.CleaningTask{},
		&models.FoodCategory{},
		&models.Food{},
		&models.Orders{},
		&models.RoomImages{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration successful!")

	return db, nil
}

// func InitDatabase() (*gorm.DB, error) {

// 	host := os.Getenv("dbHost")
// 	port := os.Getenv("dbPort")
// 	user := os.Getenv("dbUser")
// 	password := os.Getenv("dbPassword")
// 	dbname := os.Getenv("dbName")

// 	dsn := fmt.Sprintf(
// 		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
// 		host, port, user, password, dbname,
// 	)

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	if err != nil {
// 		log.Println("Database connection failed:", err)
// 		return nil, err
// 	}

// 	// Connection pool
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		log.Println("Failed to get DB instance:", err)
// 		return nil, err
// 	}
// 	sqlDB.SetMaxIdleConns(10)
// 	sqlDB.SetMaxOpenConns(100)
// 	sqlDB.SetConnMaxLifetime(time.Hour)

// 	err = db.AutoMigrate(
// 		&models.Role{},
// 		&models.User{},
// 		&models.Access{},
// 		&models.RoleAccess{},
// 		&models.UserRole{},
// 		&models.Facility{},
// 		&models.Service{},
// 		&models.Floor{},
// 		&models.RoomType{},
// 		&models.Room{},
// 		&models.Amenity{},
// 		&models.RoomAmenity{},
// 		&models.Book{},
// 		&models.BookingGuest{},
// 		&models.HistoryLog{},
// 		&models.Cleaner{},
// 		&models.CleaningTask{},
// 		&models.FoodCategory{},
// 		&models.Food{},
// 		&models.Orders{},
// 		&models.RoomImages{},
// 	)
// 	if err != nil {
// 		log.Fatal("Migration failed:", err)
// 	}

// 	log.Println("Migration successful!")

// 	return db, nil
// }
