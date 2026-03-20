package controllers

import (
	"HMS-GO/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Server) GetFoodCategory(ctx *gin.Context) {

	var foodCategory []models.FoodCategory

	if err := s.Db.Find(&foodCategory).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, foodCategory)

}

// Perform the create and update method
func (s *Server) Upsert(ctx *gin.Context) {

	var payload models.FoodCategory

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id := ctx.Param("id")
	var action string
	if id != "" {
		action = "updated"
		if err := s.Db.Model(&payload).
			Where("food_category_id", id).
			Updates(payload).Error; err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update"})
			return
		}

	} else {
		action = "inserted"
		id, err := GenerateCategoryId(s.Db)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate category Id"})
			return
		}

		payload.FoodCategoryId = id
		if err := s.Db.Create(payload).Error; err != nil {

			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to insert "})
			return
		}

	}
	ctx.JSON(http.StatusOK, gin.H{"success": fmt.Sprintf("Category %s successfully", action)})

}

// Delete food category
func (s *Server) DeleteFoodCategory(ctx *gin.Context) {
	categoryId := ctx.Param("id")

	result := s.Db.
		Where("food_category_id = ?", categoryId).
		Delete(&models.FoodCategory{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Category not found"})
		return
	}

	userId := s.GetUserId(ctx)
	err := s.CreateLogs("Food", categoryId, "Delete", "Deleted a role", "", "", userId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": "Category deleted successfully"})
}

// Generate Unique Id
func GenerateCategoryId(db *gorm.DB) (string, error) {
	var lastID string

	err := db.
		Model(&models.FoodCategory{}).
		Select("food_category_id").
		Order("food_category_id DESC").
		Limit(1).
		Scan(&lastID).Error

	if err != nil {
		return "", err
	}

	nextNumber := 1

	if lastID != "" {
		fmt.Sscanf(lastID, "CATEGORY-%d", &nextNumber)
		nextNumber++
	}

	return fmt.Sprintf("CATEGORY-%03d", nextNumber), nil
}
