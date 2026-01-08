package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// Create room aminity
func (s *Server) CreateRoomAminity(ctx *gin.Context) {

	var roomAminity models.RoomAminity
	//Validate first if
	if err := ctx.ShouldBind(&roomAminity); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Db.Create(&roomAminity).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Room aminity name already exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create Aminity",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Room aminity created successfully"})

}

// Update room aminity
func (s *Server) UpdateRoomAminity(ctx *gin.Context) {
	roomAminityId := ctx.Param("roomid")

	var payload models.RoomAminity
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Aminity{}).
		Where("room_id = ?", roomAminityId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Room Aminity updated successfully"})
}

// Delete room aminity
func (s *Server) DeleteRoomAminity(ctx *gin.Context) {
	roomAminityId := ctx.Param("roomid")

	result := s.Db.
		Where("room_id= ?", roomAminityId).
		Delete(&models.RoomAminity{})

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
func (s *Server) GetRoomAminities(ctx *gin.Context) {

	var roomAminities []models.RoomAminity

	if err := s.Db.Preload("Room").
		Preload("Amenity").
		Find(&roomAminities).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, roomAminities)

}
