package controllers

import (
	"HMS-GO/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RoomDto struct {
	RoomId     string                  `form:"roomid"`
	RoomNumber string                  `form:"roomnumber"`
	RoomTypeId string                  `form:"roomtypeid"`
	FloorId    string                  `form:"floorid"`
	Capacity   string                  `form:"capacity"`
	Price      decimal.Decimal         `form:"price"`
	Image      []*multipart.FileHeader `form:"roomImages" binding:"required"`
	Status     string                  `form:"status"`
	CreatedAt  time.Time               `form:"created_at"`
}

// Create Room
func (s *Server) CreateRoom(ctx *gin.Context) {

	var dto RoomDto

	// Validate form fields
	if err := ctx.ShouldBind(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get uploaded files manually (ShouldBind may not catch multipart slice)
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["roomImages"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one image is required"})
		return
	}

	roomID, err := GenerateRoomID(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Room ID"})
		return
	}

	userId := s.GetUserId(ctx)

	// Build room model (without images first)
	room := models.Room{
		RoomId:     roomID,
		RoomNumber: dto.RoomNumber,
		RoomTypeId: dto.RoomTypeId,
		FloorId:    dto.FloorId,
		Capacity:   dto.Capacity,
		Price:      dto.Price,
		Status:     dto.Status,
	}

	// Save room to DB
	if err := s.Db.Create(&room).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Room number already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, file := range files {
		// Sanitize filename
		ext := filepath.Ext(file.Filename)
		nameWithoutExt := strings.TrimSuffix(file.Filename, ext)
		filename := fmt.Sprintf("%s_%s%s", nameWithoutExt, roomID, ext)
		savePath := filepath.Join("src", "room_images", filename)

		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + filename})
			return
		}

		// Save image record to DB
		roomImage := models.RoomImages{
			ImageId: uuid.New().String(),
			RoomId:  roomID,
			Image:   filename,
		}
		if err := s.Db.Create(&roomImage).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image record"})
			return
		}
	}

	// Log
	if err := s.CreateLogs("Room", room.RoomId, "Create", "Created a room", "", "", userId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Room created successfully"})
}

// Update room
func (s *Server) UpdateRoom(ctx *gin.Context) {
	roomid := ctx.Param("roomid")

	userId := s.GetUserId(ctx)
	if userId == "" {
		return
	}

	// Get old room data before update
	var oldRoom models.Room
	if err := s.Db.
		Preload("RoomImages").
		Where("room_id = ?", roomid).First(&oldRoom).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room"})
		return
	}

	var dto RoomDto
	if err := ctx.ShouldBind(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build updated room payload from DTO
	payload := models.Room{
		RoomNumber: dto.RoomNumber,
		RoomTypeId: dto.RoomTypeId,
		FloorId:    dto.FloorId,
		Capacity:   dto.Capacity,
		Price:      dto.Price,
		Status:     dto.Status,
	}

	// Perform update
	if err := s.Db.Model(&models.Room{}).
		Where("room_id = ?", roomid).
		Updates(payload).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	form, err := ctx.MultipartForm()
	if err == nil {
		files := form.File["roomImages"]

		for _, file := range files {
			ext := filepath.Ext(file.Filename)
			nameWithoutExt := strings.TrimSuffix(file.Filename, ext)
			filename := fmt.Sprintf("%s_%s%s", nameWithoutExt, roomid, ext)
			savePath := filepath.Join("src", "room_images", filename)

			var existingImage models.RoomImages
			if err := s.Db.Where("room_id = ? AND image = ?", roomid, filename).First(&existingImage).Error; err == nil {
				continue
			}

			if err := ctx.SaveUploadedFile(file, savePath); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + filename})
				return
			}

			roomImage := models.RoomImages{
				ImageId: uuid.New().String(),
				RoomId:  roomid,
				Image:   filename,
			}
			if err := s.Db.Create(&roomImage).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image record"})
				return
			}
		}
	}

	// Get updated room data for logging
	var newRoom models.Room
	if err := s.Db.Where("room_id = ?", roomid).First(&newRoom).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated room"})
		return
	}

	// Log old and new values
	oldValueJSON, _ := json.Marshal(oldRoom)
	newValueJSON, _ := json.Marshal(newRoom)

	if err := s.CreateLogs(
		"Room",
		roomid,
		"UPDATE",
		"Updated room",
		string(oldValueJSON),
		string(newValueJSON),
		userId,
	); err != nil {
		fmt.Println("Failed to create log:", err)
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

	//Perfrom insert logs
	userId := s.GetUserId(ctx)
	err := s.CreateLogs("Room", roomId, "Delete", "Deleted a room", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		Preload("RoomImages").
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

// Delete role
func (s *Server) DeleteRoomImage(ctx *gin.Context) {
	filename := ctx.Param("filename")

	fmt.Println(filename)
	result := s.Db.
		Where("image = ?", filename).
		Delete(&models.RoomImages{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Room not found"})
		return
	}

	//Perfrom insert logs
	userId := s.GetUserId(ctx)
	err := s.CreateLogs("Room", filename, "Delete", "Deleted a room image", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(204)
}
