package controllers

import (
	"gorm.io/gorm"
)

// Server struct to hold the database connection
type Server struct {
	Db *gorm.DB
}

// For dependency injection itself
func NewServer(db *gorm.DB) *Server {
	return &Server{
		Db: db,
	}
}
