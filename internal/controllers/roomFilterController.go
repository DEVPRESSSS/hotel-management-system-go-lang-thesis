package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) FilterRoom(ctx *gin.Context) {

	filterMessage := ctx.Param("filter")

	if filterMessage == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "filter message must not be empty"})
		return
	}

	var rooms []models.Room
	query := s.Db.Model(&models.Room{}).
		Preload("RoomType").
		Preload("Floor").
		Joins("JOIN room_types ON room_types.room_typeid = rooms.room_type_id")

	switch filterMessage {
	case "available", "occupied":
		query = query.Where("rooms.status = ?", filterMessage)
	case "standard", "deluxe":
		query = query.Where("room_types.room_type_name = ?", filterMessage)
	case "1-2 guests":
		query = query.Where("rooms.capacity = ?", "2")
	case "3-4 guests":
		query = query.Where("rooms.capacity = ?", "4")
	}

	if err := query.Find(&rooms).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	ctx.JSON(http.StatusOK, rooms)
}
