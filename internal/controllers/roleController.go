package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Role struct {
	RoleId   string `json:"roleid"`
	RoleName string `json:"name"`
}

func (s *Server) CreateRole(ctx *gin.Context) {

	var role models.Role
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
				"error": "Username or email already exists",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Role created successfully"})

}

func (s *Server) GetRoles(ctx *gin.Context) {

	var roles []models.Role

	if err := s.Db.Find(&roles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, roles)

}

// Fetch the information of the selected record
func (s *Server) GetRole(ctx *gin.Context) {

	roleID := ctx.Param("roleid")

	var role models.Role
	if err := s.Db.
		Where("role_id = ?", roleID).First(&role).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": role})
}
