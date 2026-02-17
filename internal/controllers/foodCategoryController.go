package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
