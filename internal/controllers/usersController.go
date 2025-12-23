package controllers

import (
	"HMS-GO/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

	//Validate first if there are input error
	if err := s.Db.Create(&user).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
		fmt.Println(err)
		return
	}

	//Return 200  if the input succeed
	ctx.JSON(http.StatusOK, gin.H{"Success:": "User created created successfully"})

}


