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

// Create role
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

// Update role
func (s *Server) UpdateRole(ctx *gin.Context) {
	userID := ctx.Param("roleid")

	var payload models.Role
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Role{}).
		Where("role_id = ?", userID).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Role updated successfully"})
}

// Delete role
func (s *Server) DeleteRole(ctx *gin.Context) {
	roleid := ctx.Param("roleid")

	result := s.Db.
		Where("role_id = ?", roleid).
		Delete(&models.Role{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Role not found"})
		return
	}

	ctx.Status(204)
}

// Get all the roles from db
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

// Fetch the information of the selected record in role
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
