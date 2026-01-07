package controllers

import (
	"HMS-GO/internal/models"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// Create role access
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
func (s *Server) UpdateRoleAcccess(ctx *gin.Context) {
	roleAccessID := ctx.Param("roleid")

	var payload models.RoleAccess
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

	ctx.JSON(200, gin.H{"success": "Role access updated successfully"})
}

// Delete role access
func (s *Server) DeleteRoleAccess(ctx *gin.Context) {
	roleid := ctx.Param("accessid")

	result := s.Db.
		Where("access_id = ?", roleid).
		Delete(&models.RoleAccess{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Role access not found"})
		return
	}

	ctx.Status(204)
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

// =====================================
// =======Access Management=============
// Create access function
func (s *Server) CreateAccess(ctx *gin.Context) {
	var access models.Access

	if err := ctx.ShouldBindJSON(&access); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Printf("BOUND: %+v\n", access)

		return
	}

	if err := s.Db.Create(&access).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create access",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "Access created successfully",
	})
}

// Update access
func (s *Server) UpdateAcccess(ctx *gin.Context) {
	accessId := ctx.Param("accessid")

	var payload models.Access
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.Access{}).
		Where("access_id = ?", accessId).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "Access updated successfully"})
}

// Delete access
func (s *Server) DeleteAccess(ctx *gin.Context) {
	roleid := ctx.Param("accessid")

	result := s.Db.
		Where("access_id = ?", roleid).
		Delete(&models.Access{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Role access not found"})
		return
	}

	ctx.Status(204)
}

// Fetch the information of the selected record in role
func (s *Server) GetAccess(ctx *gin.Context) {

	accessId := ctx.Param("accessid")

	var access models.Access
	if err := s.Db.
		Where("access_id = ?", accessId).First(&access).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": access})
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
