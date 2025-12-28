package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Facility struct {
	FacilityId      string    `json:"facility_id"`
	FacilityName    string    `json:"facility_name"`
	Status          bool      `json:"status"`
	MaintenanceDate time.Time `json:"maintenance_date"`
	CreatedAt       time.Time `json:"created_at"`
}

// Create facility
func (s *Server) CreateFacility(ctx *gin.Context) {

	var facility models.Facility
	//Validate first if
	if err := ctx.ShouldBind(&facility); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create Role error handling
	if err := s.Db.Create(&facility).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Facility already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Facility created successfully"})

}

// Update role
func (s *Server) UpdateFacility(ctx *gin.Context) {
	facilityId := ctx.Param("facilityid")

	var payload models.Facility
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Facility{}).
		Where("facility_id = ?", facilityId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Facility updated successfully"})
}

// Delete role
func (s *Server) Deletefacility(ctx *gin.Context) {
	faciltyId := ctx.Param("facilityid")

	result := s.Db.
		Where("facility_id = ?", faciltyId).
		Delete(&models.Facility{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Facility not found"})
		return
	}

	ctx.Status(204)
}

// Get all the roles from db
func (s *Server) GetFacilities(ctx *gin.Context) {

	var facilities []models.Facility

	if err := s.Db.Find(&facilities).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, facilities)

}

// Fetch the information of the selected record in role
func (s *Server) GetFacility(ctx *gin.Context) {

	facilityId := ctx.Param("facilityid")

	var facility models.Facility
	if err := s.Db.
		Where("facility_id = ?", facilityId).First(&facility).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": facility})
}
