package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskflow/config"
	"taskflow/models"
)

// GetNotifications GET /notifications
// @Summary      Get User Notification
// @Tags         Notifications
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200 {array} models.Notification
// @Router       /notifications [get]
func GetNotifications(c *gin.Context) {
	userID, _ := c.Get("userID")
	rows, _ := config.DB.Query(
		"SELECT id, user_id, task_id, message, is_read, created_at FROM notifications WHERE user_id = ? AND is_read = FALSE",
		userID,
	)
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		rows.Scan(&n.ID, &n.UserID, &n.TaskID, &n.Message, &n.IsRead, &n.CreatedAt)
		notifs = append(notifs, n)
	}
	c.JSON(http.StatusOK, notifs)
}
