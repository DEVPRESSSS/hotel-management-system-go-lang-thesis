package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Create user
func (s *Server) RegisterGuest(ctx *gin.Context) {

	var create models.RegisterInput

	if err := ctx.ShouldBind(&create); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleId, err := s.GetGuestRoleID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "role id is not found"})
		return
	}

	guestID, err := GenerateGuestID(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Amenity ID",
		})
		return
	}

	create.UserId = guestID
	create.RoleId = roleId
	//Assign the value of each input
	user := models.User{
		UserId:   create.UserId,
		Username: create.Username,
		FullName: create.FullName,
		Email:    create.Email,
		Password: create.Password,
		RoleId:   create.RoleId,
		Locked:   create.Locked,
	}
	utils.HashPassword(&user)

	if err := s.Db.Create(&user).Error; err != nil {

		//Check first if the error number is 1062
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Username is already taken",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Registration failed",
		})
		return
	}
	//Return 200  if the input succeed
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Registration succeed",
	})
}

// Get all the roles from db
func (s *Server) GetGuestRoleID() (string, error) {
	var role models.Role

	err := s.Db.
		Where("role_name = ?", "Guest").
		First(&role).Error

	if err != nil {
		return "", err
	}

	return role.RoleId, nil
}

// Generate auto IncrementId
func GenerateGuestID(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.User{}).
		Select("user_id").
		Order("user_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "GUEST-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("GUEST-%03d", nextNumber), nil
}
