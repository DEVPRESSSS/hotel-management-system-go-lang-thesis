package controllers

import (
	"HMS-GO/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Role struct {
	RoleId   string `json:"roleid"`
	RoleName string `json:"name"`
}

func (s *Server) GetRoles(ctx *gin.Context) {

	var roles []models.Role

	if err := s.Db.Find(&roles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Print(roles)
	ctx.JSON(http.StatusOK, roles)

}
