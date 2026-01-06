package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

func (s *Server) CreateRoleAccess(ctx *gin.Context) {

	var roleAccess models.RoleAccess

	// Bind request
	if err := ctx.ShouldBind(&roleAccess); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create record
	if err := s.Db.Create(&roleAccess).Error; err != nil {

		// Handle duplicate entry (MySQL 1062)
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "Role access already exists",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create role access",
		})
		return
	}

	// Success
	ctx.JSON(http.StatusOK, gin.H{
		"success": "Role access created successfully",
	})
}

// Update role access
func (s *Server) UpdateAcccess(ctx *gin.Context) {
	roleAccessID := ctx.Param("roleid")

	var payload models.Role
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Role{}).
		Where("role_id = ?", roleAccessID).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Role updated successfully"})
}

// Create access
func (s *Server) CreateAccess(ctx *gin.Context) {

	var access models.Access
	if err := ctx.ShouldBind(&access); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Create access error handling
	if err := s.Db.Create(&access).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Accesss created successfully"})
}

// fetch all the role access properly
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

// Fetch all the access
func (s *Server) Access(ctx *gin.Context) {

	var access []models.Access

	if err := s.Db.
		Find(&access).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, access)
}
