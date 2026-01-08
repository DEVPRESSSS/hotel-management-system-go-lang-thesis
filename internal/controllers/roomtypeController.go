package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func (s *Server) GetRoomtype(ctx *gin.Context) {

	var rt []models.RoomType

	if err := s.Db.Find(&rt).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, rt)

}

// Create room type
func (s *Server) CreateRoomType(ctx *gin.Context) {

	var role models.RoomType
	//Validate first if
	if err := ctx.ShouldBind(&role); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create Role error handling
	if err := s.Db.Create(&role).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Room type name already exist",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Role created successfully"})

}

// Update room type
func (s *Server) UpdateRoomType(ctx *gin.Context) {
	roomTypeId := ctx.Param("roomtypeid")

	var payload models.RoomType
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.RoomType{}).
		Where("room_typeid = ?", roomTypeId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Room type updated successfully"})
}

// Delete room type
func (s *Server) DeleteRoomType(ctx *gin.Context) {
	roleid := ctx.Param("roomtypeid")

	result := s.Db.
		Where("room_typeid = ?", roleid).
		Delete(&models.RoomType{})

	if result.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Room type name already exist",
			})
			return
		} else if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1451 {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "This room type has related to other records, deletion failed!!",
			})
			return
		}
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Role not found"})
		return
	}

	ctx.Status(204)
}

// Fetch the information of the selected record in role
func (s *Server) GetRoomTypeRecord(ctx *gin.Context) {

	roomtypeid := ctx.Param("roomtypeid")

	var rt models.RoomType
	if err := s.Db.
		Where("room_typeid = ?", roomtypeid).First(&rt).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": rt})
}
