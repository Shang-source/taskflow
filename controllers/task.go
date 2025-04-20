package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"taskflow/config"
	"taskflow/models"
)

// GetTasks GET /tasks
// GetTasks GET /tasks
// @Summary      Get all tasks of the current user
// @Description  Supports pagination and sorting (if needed)
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200 {array}  models.Task
// @Router       /tasks [get]
func GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID") // JWT user ID
	rows, err := config.DB.Query("SELECT id, title, description, is_done, user_id, due_date FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.IsDone, &t.UserID, &t.DueDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tasks = append(tasks, t)
	}
	c.JSON(http.StatusOK, tasks)
}
