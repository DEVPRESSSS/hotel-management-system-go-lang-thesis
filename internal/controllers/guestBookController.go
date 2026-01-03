package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RoomSelected(ctx *gin.Context) {
	roomId := ctx.Param("roomid")
	var room models.Room
	if err := s.Db.
		Where("room_id = ?", &roomId).
		Find(&room).Error; err != nil {

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
}
