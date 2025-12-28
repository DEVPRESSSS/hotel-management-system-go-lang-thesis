package utils

import (
	"HMS-GO/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(user *models.User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		return err
	}

	user.Password = string(hashedPassword)
	return nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
