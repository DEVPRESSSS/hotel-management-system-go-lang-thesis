package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type FoodRequest struct {
	Name           string                `form:"name" binding:"required"`
	Description    string                `form:"description" binding:"required"`
	Image          *multipart.FileHeader `form:"image" binding:"required"`
	FoodCategoryId string                `form:"food-category-id" binding:"required"`
	Price          decimal.Decimal       `form:"price" binding:"required"`
	Status         string                `form:"status" binding:"required"`
}

// Create food service
func (s *Server) CreateFoodService(ctx *gin.Context) {

	var req FoodRequest
	//Validate first
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foodId, err := GenerateFoodServiceId(s.Db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate food service",
		})
		return
	}

	ext := filepath.Ext(req.Image.Filename)
	nameWithoutExt := strings.TrimSuffix(req.Image.Filename, ext)
	filename := fmt.Sprintf("%s_%s%s", nameWithoutExt, foodId, ext)
	savePath := filepath.Join("src", "food_images", filename)

	//Production code
	err = ctx.SaveUploadedFile(req.Image, savePath)

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, "Unknown error")
		return
	}
	food := models.Food{
		FoodId:         foodId,
		Name:           req.Name,
		Description:    req.Description,
		Image:          filename,
		FoodCategoryId: req.FoodCategoryId,
		Price:          req.Price,
		Status:         req.Status,
	}
	if err := s.Db.
		Create(&food).Error; err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Food service already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
	}

	userId := s.GetUserId(ctx)
	err = s.CreateLogs("Food Service", food.FoodId, "Create", "Created a food service", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Food service created successfully"})

}

// Update food service
func (s *Server) UpdateFoodService(ctx *gin.Context) {
	foodServiceId := ctx.Param("id")

	var payload models.Food
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Food{}).
		Where("food_id = ?", foodServiceId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	userId := s.GetUserId(ctx)
	err := s.CreateLogs("Food service", foodServiceId, "Update", "Updated a food service", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"success": "Food service updated successfully"})
}

// Delete role
func (s *Server) DeleteFoodService(ctx *gin.Context) {
	foodId := ctx.Param("id")

	var obj models.Food
	if err := s.Db.Where("food_id = ?", foodId).First(&obj).Error; err != nil {

		ctx.JSON(http.StatusBadRequest, "bad request")
		return
	}
	imagePath := obj.Image
	if imagePath == "" {
		fmt.Print("Failed to fetch image")
		return
	}

	//Production code
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to fetch wd")
		return
	}

	path := filepath.Join(cwd, "src", "food_images", imagePath)
	os.Remove(path)

	//os.Remove(imagePath)

	result := s.Db.
		Where("food_id = ?", foodId).
		Delete(&models.Food{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Food not found"})
		return
	}

	userId := s.GetUserId(ctx)
	err = s.CreateLogs("Food Service", foodId, "Delete", "Deleted a food service", "", "", userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(204)
}

// Get all the food services from db
func (s *Server) GetFoodServices(ctx *gin.Context) {

	var foodServices []models.Food

	if err := s.Db.
		Preload("FoodCategory").
		Find(&foodServices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, foodServices)

}

// Fetch the information of the selected record in food service
func (s *Server) GetFoodService(ctx *gin.Context) {

	foodId := ctx.Param("id")

	var service models.Food
	if err := s.Db.
		Where("food_id = ?", foodId).First(&service).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": service})
}

// Generate auto IncrementId
func GenerateFoodServiceId(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.Food{}).
		Select("food_id").
		Order("food_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "FOOD-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("FOOD-%03d", nextNumber), nil
}
