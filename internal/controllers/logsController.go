package controllers

import (
	"HMS-GO/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get all the roles from db
func (s *Server) GetLogs(ctx *gin.Context) {

	var logs []models.HistoryLog

	if err := s.Db.Find(&logs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, logs)

}

func (s *Server) CreateLogs(EntityType string, EntityID string, Action string, Description string, OldValue string, NewValue string, UserId string) error {

	logId, err := GenerateLogId(s.Db)
	if err != nil {

		return err
	}
	logs := models.HistoryLog{
		Id:          logId,
		EntityType:  EntityType,
		EntityID:    EntityID,
		Action:      Action,
		Description: Description,
		OldValue:    OldValue,
		NewValue:    NewValue,
		PerformedBy: UserId,
	}
	if err := s.Db.Create(&logs).Error; err != nil {

		return err
	}

	return nil
}
