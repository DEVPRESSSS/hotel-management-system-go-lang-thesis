package controllers

import (
	"HMS-GO/internal/models"
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

func (s *Server) CreateUser(ctx *gin.Context) {

	var create models.CreateUserInput

	if err := ctx.ShouldBind(&create); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(create)

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
