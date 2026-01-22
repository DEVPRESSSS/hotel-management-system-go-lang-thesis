package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Create aminity
func (s *Server) CreateAminity(ctx *gin.Context) {
	var amenity models.Amenity

	// Bind request
	if err := ctx.ShouldBind(&amenity); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate sequential ID
	amenityID, err := GenerateAmenityID(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Amenity ID",
		})
		return
	}

	amenity.AmenityId = amenityID

	// Create record
	if err := s.Db.Create(&amenity).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Amenity name already exists",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create Amenity",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":    "Amenity created successfully",
		"amenity_id": amenity.AmenityId,
	})
}

// Update aminity
func (s *Server) UpdateAminity(ctx *gin.Context) {
	aminityId := ctx.Param("amenityid")

	var payload models.Amenity
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Amenity{}).
		Where("amenity_id = ?", aminityId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Aminity updated successfully"})
}

// Delete aminity
func (s *Server) DeleteAminity(ctx *gin.Context) {
	aminityId := ctx.Param("amenityid")

	result := s.Db.
		Where("amenity_id= ?", aminityId).
		Delete(&models.Amenity{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Amenity already exist"})
		return
	}

	ctx.Status(204)
}

// Fetch the information of the selected record in aminity
func (s *Server) GetAminity(ctx *gin.Context) {

	aminityId := ctx.Param("amenityid")

	var aminity models.Amenity
	if err := s.Db.
		Where("amenity_id = ?", aminityId).First(&aminity).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": aminity})
}

// Get all the aminity from db
func (s *Server) GetAminities(ctx *gin.Context) {

	var aminities []models.Amenity

	if err := s.Db.Find(&aminities).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, aminities)

}

// Generate auto IncrementId
func GenerateAmenityID(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.Amenity{}).
		Select("amenity_id").
		Order("amenity_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "AMENITY-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("AMENITY-%03d", nextNumber), nil
}
