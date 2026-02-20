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
	// Assign all 4 accesses to Admin role
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

	if err := SeedAdminUser(db); err != nil {
		return err
	}

	log.Println("Seeders completed!")
	return nil
}
