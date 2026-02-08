package controllers

import (
	"HMS-GO/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Room struct {
	RoomId     string          `json:"roomid"`
	RoomNumber string          `json:"roomnumber"`
	RoomTypeId string          `json:"roomtypeid"`
	FloorId    string          `json:"floorid"`
	Capacity   string          `json:"capacity"`
	Price      decimal.Decimal `json:"price"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
}

// Create Room
func (s *Server) CreateRoom(ctx *gin.Context) {

	var room models.Room
	//Validate first
	if err := ctx.ShouldBind(&room); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID, err := GenerateRoomID(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Amenity ID",
		})
		return
	}

	userId := s.GetUserId(ctx)

	//Asign the room id here
	room.RoomId = roomID
	//Perform insert
	if err := s.Db.Create(&room).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Room number already exists",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Perfrom insert logs
	err = s.CreateLogs("Room", room.RoomId, "Create", "Created a room", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "Room created successfully"})

}

// Update room
func (s *Server) UpdateRoom(ctx *gin.Context) {
	roomid := ctx.Param("roomid")

	// Get user ID first
	userId := s.GetUserId(ctx)
	if userId == "" {
		return
	}

	// Get old room data before update
	var oldRoom models.Room
	if err := s.Db.Where("room_id = ?", roomid).First(&oldRoom).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room"})
		return
	}

	// Bind new data
	var payload models.Room
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Perform update
	if err := s.Db.Model(&models.Room{}).
		Where("room_id = ?", roomid).
		Updates(payload).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	// Get updated room data
	var newRoom models.Room
	if err := s.Db.Where("room_id = ?", roomid).First(&newRoom).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated room"})
		return
	}

	// Convert to JSON strings for logging
	oldValueJSON, _ := json.Marshal(oldRoom)
	newValueJSON, _ := json.Marshal(newRoom)

	// Create log with old and new values
	err := s.CreateLogs(
		"Room",
		roomid,
		"UPDATE",
		"Updated room",
		string(oldValueJSON),
		string(newValueJSON),
		userId,
	)
	if err != nil {
		fmt.Println("Failed to create log:", err)
		// Don't fail the update if logging fails
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Room updated successfully",
		"data":    newRoom,
	})
}

// Delete role
func (s *Server) DeleteRoom(ctx *gin.Context) {
	roomId := ctx.Param("roomid")

	result := s.Db.
		Where("room_id = ?", roomId).
		Delete(&models.Room{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Room not found"})
		return
	}

	ctx.Status(204)
}

// Get all the services from db
func (s *Server) GetRooms(ctx *gin.Context) {

	var rooms []models.Room

	if err := s.Db.
		Preload("Floor").
		Preload("RoomType").
		Find(&rooms).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, rooms)

}

// Fetch the information of the selected record in role
func (s *Server) GetRoom(ctx *gin.Context) {

	roomId := ctx.Param("roomid")

	var room models.Room
	if err := s.Db.
		Where("room_id = ?", roomId).First(&room).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": room})
}

// Generate auto IncrementId
func GenerateRoomID(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.Room{}).
		Select("room_id").
		Order("room_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "ROOM-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("ROOM-%03d", nextNumber), nil
}

// Generate auto IncrementId
func GenerateLogId(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.HistoryLog{}).
		Select("history_id").
		Order("history_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "LOG-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("LOG-%03d", nextNumber), nil
}

func (s *Server) CreateLogs(EntityType string, EntityID string, Action string, Description string, OldValue string, NewValue string, UserId string) error {

	logId, err := GenerateLogId(s.Db)
	if err != nil {

		return err
	}
	logs := models.HistoryLog{
		Id:          logId,
		EntityType:  EntityType,
		EntityID:    EntityID,
		Action:      Action,
		Description: Description,
		OldValue:    OldValue,
		NewValue:    NewValue,
		PerformedBy: UserId,
	}
	if err := s.Db.Create(&logs).Error; err != nil {

		return err
	}

	return nil
}
