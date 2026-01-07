package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// Create aminity
func (s *Server) CreateAminity(ctx *gin.Context) {

	var aminity models.Aminity
	//Validate first if
	if err := ctx.ShouldBind(&aminity); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create aminity error handling
	if err := s.Db.Create(&aminity).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Aminity name already exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create Aminity",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Aminity created successfully"})

}

// Update aminity
func (s *Server) UpdateAminity(ctx *gin.Context) {
	aminityId := ctx.Param("aminityid")

	var payload models.Aminity
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Aminity{}).
		Where("aminity_id = ?", aminityId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Aminity updated successfully"})
}

// Delete aminity
func (s *Server) DeleteAminity(ctx *gin.Context) {
	aminityId := ctx.Param("aminityid")

	result := s.Db.
		Where("aminity_id= ?", aminityId).
		Delete(&models.Aminity{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Role not found"})
		return
	}

	ctx.Status(204)
}

// Get all the aminity from db
func (s *Server) GetAminities(ctx *gin.Context) {

	var aminities []models.Aminity

	if err := s.Db.Find(&aminities).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, aminities)

}
