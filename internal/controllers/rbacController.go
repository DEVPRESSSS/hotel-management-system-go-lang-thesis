package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RoleAccess(ctx *gin.Context) {

	var roleAccess []models.RoleAccess

	if err := s.Db.
		Preload("Access").
		Preload("Role").
		Find(&roleAccess).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, roleAccess)
}
