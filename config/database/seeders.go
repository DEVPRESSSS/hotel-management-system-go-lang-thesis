package database

import (
	"HMS-GO/internal/models"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{RoleId: "role-admin", RoleName: "Admin"},
		{RoleId: "role-frontdesk", RoleName: "Frontdesk"},
		{RoleId: "role-guest", RoleName: "Guest"},
	}

	for _, role := range roles {
		var existing models.Role
		result := db.Where("role_id = ?", role.RoleId).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to seed role '%s': %v", role.RoleName, err)
				return err
			}
			log.Printf("Seeded role: %s", role.RoleName)
		} else {
			log.Printf("Role '%s' already exists, skipping", role.RoleName)
		}
	}
	return nil
}

func SeedAccess(db *gorm.DB) error {
	accesses := []models.Access{
		{AccessId: "access-create", AccessName: "create"},
		{AccessId: "access-read", AccessName: "read"},
		{AccessId: "access-update", AccessName: "update"},
		{AccessId: "access-delete", AccessName: "delete"},
	}

	for _, access := range accesses {
		var existing models.Access
		result := db.Where("access_id = ?", access.AccessId).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&access).Error; err != nil {
				log.Printf("Failed to seed access '%s': %v", access.AccessName, err)
				return err
			}
			log.Printf("Seeded access: %s", access.AccessName)
		} else {
			log.Printf("Access '%s' already exists, skipping", access.AccessName)
		}
	}
	return nil
}

func SeedRoleAccess(db *gorm.DB) error {
	adminAccesses := []models.RoleAccess{
		{RoleID: "role-admin", AccessID: "access-create"},
		{RoleID: "role-admin", AccessID: "access-read"},
		{RoleID: "role-admin", AccessID: "access-update"},
		{RoleID: "role-admin", AccessID: "access-delete"},
	}

	for _, ra := range adminAccesses {
		var existing models.RoleAccess
		result := db.Where("role_id = ? AND access_id = ?", ra.RoleID, ra.AccessID).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&ra).Error; err != nil {
				log.Printf("Failed to seed role_access role:%s access:%s — %v", ra.RoleID, ra.AccessID, err)
				return err
			}
			log.Printf("Seeded role_access: %s -> %s", ra.RoleID, ra.AccessID)
		} else {
			log.Printf("Role_access %s -> %s already exists, skipping", ra.RoleID, ra.AccessID)
		}
	}
	return nil
}

func SeedFloors(db *gorm.DB) error {
	floors := []models.Floor{
		{FloorId: "FLOOR-01", FloorName: "Floor 1"},
		{FloorId: "FLOOR-02", FloorName: "Floor 2"},
		{FloorId: "FLOOR-3", FloorName: "Floor 3"},
		{FloorId: "FLOOR-04", FloorName: "Floor 4"},
		{FloorId: "FLOOR-05", FloorName: "Floor 5"},
	}

	for _, floor := range floors {
		var existing models.Floor
		result := db.Where("floor_id = ?", floor.FloorId).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&floor).Error; err != nil {
				log.Printf("Failed to seed floor '%s': %v", floor.FloorName, err)
				return err
			}
			log.Printf("Seeded floor: %s", floor.FloorName)
		} else {
			log.Printf("Floor '%s' already exists, skipping", floor.FloorName)
		}
	}
	return nil
}

func SeedRoomTypes(db *gorm.DB) error {
	roomTypes := []models.RoomType{
		{
			RoomTypeId:   "room-001",
			RoomTypeName: "Standard",
			Description:  "A comfortable standard room with essential amenities.",
		},
		{
			RoomTypeId:   "room-002",
			RoomTypeName: "Deluxe",
			Description:  "A spacious deluxe room with premium furnishings and amenities.",
		},
		{
			RoomTypeId:   "room-003",
			RoomTypeName: "Suite",
			Description:  "A luxurious suite with a separate living area and top-tier amenities.",
		},
		{
			RoomTypeId:   "room-004",
			RoomTypeName: "Family",
			Description:  "A large room designed to accommodate families comfortably.",
		},
		{
			RoomTypeId:   "room-005",
			RoomTypeName: "Penthouse",
			Description:  "An exclusive penthouse with panoramic views and premium facilities.",
		},
	}

	for _, rt := range roomTypes {
		var existing models.RoomType
		result := db.Where("room_typeid = ?", rt.RoomTypeId).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&rt).Error; err != nil {
				log.Printf("Failed to seed room type '%s': %v", rt.RoomTypeName, err)
				return err
			}
			log.Printf("Seeded room type: %s", rt.RoomTypeName)
		} else {
			log.Printf("Room type '%s' already exists, skipping", rt.RoomTypeName)
		}
	}
	return nil
}

func SeedFoodCategories(db *gorm.DB) error {
	categories := []models.FoodCategory{
		{FoodCategoryId: "foodcat-001", Name: "Breakfast", Time: "6:00 AM - 10:00 AM"},
		{FoodCategoryId: "foodcat-002", Name: "Lunch", Time: "11:00 AM - 2:00 PM"},
		{FoodCategoryId: "foodcat-003", Name: "Dinner", Time: "6:00 PM - 10:00 PM"},
		{FoodCategoryId: "foodcat-004", Name: "Snacks", Time: "2:00 PM - 5:00 PM"},
		{FoodCategoryId: "foodcat-005", Name: "Drinks", Time: "All Day"},
	}

	for _, cat := range categories {
		var existing models.FoodCategory
		result := db.Where("food_category_id = ?", cat.FoodCategoryId).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&cat).Error; err != nil {
				log.Printf("Failed to seed food category '%s': %v", cat.Name, err)
				return err
			}
			log.Printf("Seeded food category: %s", cat.Name)
		} else {
			log.Printf("Food category '%s' already exists, skipping", cat.Name)
		}
	}
	return nil
}

func SeedAdminUser(db *gorm.DB) error {
	var existing models.User

	email := os.Getenv("superAdmin")
	password := os.Getenv("password")

	result := db.Where("email = ?", email).First(&existing)
	if result.Error != gorm.ErrRecordNotFound {
		log.Println("Admin user already exists, skipping")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return err
	}

	admin := models.User{
		UserId:    "user-admin-001",
		Username:  "xmontemorjerald",
		FullName:  "Jerald Montemor",
		Email:     "xmontemorjerald@gmail.com",
		Password:  string(hashedPassword),
		Verified:  true,
		Locked:    false,
		CreatedAt: time.Now(),
		RoleId:    "role-admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to seed admin user: %v", err)
		return err
	}

	log.Println("Seeded admin user: xmontemorjerald@gmail.com")
	return nil
}

func RunSeeders(db *gorm.DB) error {
	log.Println("Running seeders...")

	if err := SeedRoles(db); err != nil {
		return err
	}
	if err := SeedAccess(db); err != nil {
		return err
	}
	if err := SeedRoleAccess(db); err != nil {
		return err
	}
	if err := SeedFloors(db); err != nil {
		return err
	}
	if err := SeedRoomTypes(db); err != nil {
		return err
	}
	if err := SeedFoodCategories(db); err != nil {
		return err
	}
	if err := SeedAdminUser(db); err != nil {
		return err
	}

	log.Println("Seeders completed!")
	return nil
}
