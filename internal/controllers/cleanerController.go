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

// Create maintenance
func (s *Server) CreateCleaner(ctx *gin.Context) {

	var cleaner models.Cleaner
	//Validate first if
	if err := ctx.ShouldBind(&cleaner); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := GenerateId(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Attendant Id",
		})
		return
	}

	cleaner.Id = id
	//Create Role error handling
	if err := s.Db.Create(&cleaner).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Maintenance name already exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
	}

	userId := s.GetUserId(ctx)
	err = s.CreateLogs("Cleaner", cleaner.Id, "Create", "Created a cleaner", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Cleaner created successfully"})

}

// Update maintenance
func (s *Server) UpdateCleaner(ctx *gin.Context) {
	cleanerId := ctx.Param("id")

	var payload models.Cleaner
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Cleaner{}).
		Where("id = ?", cleanerId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	userId := s.GetUserId(ctx)
	err := s.CreateLogs("Cleaner", cleanerId, "Update", "Updated a cleaner", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"success": "Cleaner updated successfully"})
}

// Delete maintenance
func (s *Server) DeleteCleaner(ctx *gin.Context) {
	cleanerId := ctx.Param("id")

	result := s.Db.
		Where("id = ?", cleanerId).
		Delete(&models.Cleaner{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Cleaner name not found"})
		return
	}

	userId := s.GetUserId(ctx)
	err := s.CreateLogs("Cleaner", cleanerId, "Delete", "Deleted a cleaner", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(204)
}

// Get all the maintenance from db
func (s *Server) GetAllCleaners(ctx *gin.Context) {

	var cleaners []models.Cleaner

	if err := s.Db.Find(&cleaners).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, cleaners)

}

// Fetch the information of the selected record in maintenance
func (s *Server) GetCleaner(ctx *gin.Context) {

	id := ctx.Param("id")

	var attendant models.Cleaner
	if err := s.Db.
		Where("id = ?", id).First(&attendant).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": attendant})
}

func GenerateId(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.Cleaner{}).
		Select("id").
		Order("id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "ATTENDANT-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("ATTENDANT-%03d", nextNumber), nil
}
