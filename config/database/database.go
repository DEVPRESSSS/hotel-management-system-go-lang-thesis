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

	// err = db.AutoMigrate(&models.Role{})
	// err = db.AutoMigrate(&models.User{})
	err = db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Facility{},
		&models.Service{},
		&models.Floor{},
		&models.RoomType{},
		&models.Room{},
	)

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration successful!")

	return db, nil
}
