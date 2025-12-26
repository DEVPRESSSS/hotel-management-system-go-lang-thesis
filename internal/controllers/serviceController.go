package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Service struct {
	ServiceId   string    `json:"serviceid"`
	ServiceName string    `json:"servicename"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
}

// Create Service
func (s *Server) CreateService(ctx *gin.Context) {

	var service models.Service
	//Validate first
	if err := ctx.ShouldBind(&service); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create Role error handling
	if err := s.Db.Create(&service).Error; err != nil {

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

	ctx.JSON(http.StatusOK, gin.H{"success": "Service created successfully"})

}

// Update role
func (s *Server) UpdateService(ctx *gin.Context) {
	serviceId := ctx.Param("serviceid")

	var payload models.Service
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Service{}).
		Where("service_id = ?", serviceId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Facility updated successfully"})
}

// Delete role
func (s *Server) DeleteService(ctx *gin.Context) {
	serviceId := ctx.Param("serviceid")

	result := s.Db.
		Where("service_id = ?", serviceId).
		Delete(&models.Service{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Service not found"})
		return
	}

	ctx.Status(204)
}

// Get all the services from db
func (s *Server) GetServices(ctx *gin.Context) {

	var services []models.Service

	if err := s.Db.Find(&services).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, services)

}

// Fetch the information of the selected record in role
func (s *Server) GetService(ctx *gin.Context) {

	serviceId := ctx.Param("serviceid")

	var service models.Service
	if err := s.Db.
		Where("service_id = ?", serviceId).First(&service).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": service})
}
