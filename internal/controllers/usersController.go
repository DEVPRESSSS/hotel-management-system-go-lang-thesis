package controllers

import (
	"HMS-GO/internal/models"
	"HMS-GO/internal/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type User struct {
	UserId   string `json:"userid"`
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   string `json:"roleid"`
	Locked   bool   `json:"locked"`
}

// Create user
func (s *Server) CreateUser(ctx *gin.Context) {

	var create models.CreateUserInput

	if err := ctx.ShouldBind(&create); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Assign the value of each input
	user := models.User{
		UserId:   create.UserId,
		Username: create.Username,
		FullName: create.FullName,
		Email:    create.Email,
		Password: create.Password,
		RoleId:   create.RoleId,
		Locked:   create.Locked,
	}
	utils.HashPassword(&user)

	if err := s.Db.Create(&user).Error; err != nil {

		//Check first if the error number is 1062
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
		return
	}
	//Return 200  if the input succeed
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

// Delete user
func (s *Server) DeleteUser(ctx *gin.Context) {
	userid := ctx.Param("userid")

	result := s.Db.
		Where("user_id = ?", userid).
		Delete(&models.User{})

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}

	ctx.Status(204)
}

func (s *Server) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("userid")

	var payload models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}

	if err := s.Db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(payload).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Update failed"})
		return
	}

	ctx.JSON(200, gin.H{"success": "User updated successfully"})
}

// Fetch the information of the selected record
func (s *Server) GetUser(ctx *gin.Context) {

	userID := ctx.Param("userid")

	var user models.User
	if err := s.Db.
		Where("user_id = ?", userID).First(&user).Error; err != nil {
		fmt.Print(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Error fetching data!!!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": user})
}

// Fetch all the data from the database
func (s *Server) GetAllUsers(ctx *gin.Context) {

	var users []models.User

	if err := s.Db.
		Preload("Role").
		Find(&users).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// Fetch all the data from the database
func (s *Server) GetAllGuest(ctx *gin.Context) {
	var users []models.User

	err := s.Db.
		Joins("JOIN roles ON roles.role_id = users.role_id").
		Where("roles.role_name = ?", "Guest").
		Preload("Role").
		Find(&users).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
