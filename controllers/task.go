package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"taskflow/config"
	"taskflow/models"
)

// CreateTask POST /tasks
// @Summary      Create Tasks
// @Description  为当前用户创建一个新任务
// @Tags         Tasks
// @Security     ApiKeyAuth
// 若用 JWT，可加这一行
// @Accept       json
// @Produce      json
// @Param        task  body  models.Task  true  "任务内容"
// @Success      201   {object}  models.Task
// @Router       /tasks [post]
func CreateTask(c *gin.Context) {
	userID, _ := c.Get("userID") // 来自 JWT
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*result, err := config.DB.Exec("INSERT INTO tasks (title, description, user_id) VALUES (?,?,?)",
	task.Title, task.Description, userID) */
	result, err := config.DB.Exec(
		"INSERT INTO tasks (title, description, is_done, due_date, reminded, user_id) VALUES (?,?,?,?,?,?)",
		task.Title, task.Description, task.IsDone, task.DueDate, false, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	task.ID = id
	task.UserID = userID.(int64)

	c.JSON(http.StatusCreated, gin.H{"message": "任务创建成功", "task": task})
}

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
	userID, _ := c.Get("userID") // JWT 中的用户ID
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

// GetTaskByID GET /tasks/:id
// GetTaskByID GET /tasks/{id}
// @Summary      Get a single task by ID
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      int  true  "任务 ID"
// @Success      200  {object}  models.Task
// @Failure      404  {object}  models.ErrorResponse
// @Router       /tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var t models.Task
	err = config.DB.QueryRow("SELECT id, title, description, is_done, user_id, due_date FROM tasks WHERE id = ? AND user_id = ?",
		id, userID).Scan(&t.ID, &t.Title, &t.Description, &t.IsDone, &t.UserID, &t.DueDate)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或不属于当前用户"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

// UpdateTask PUT /tasks/:id
// UpdateTask PUT /tasks/{id}
// @Summary      Update Tasks
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "任务 ID"
// @Param        task  body      models.Task  true  "任务内容"
// @Success      200   {object}  models.Task
// @Router       /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*_, err = config.DB.Exec("UPDATE tasks SET title=?, description=?, is_done=? WHERE id=? AND user_id=?",
	input.Title, input.Description, input.IsDone, id, userID)   */
	_, err = config.DB.Exec(
		`UPDATE tasks 
                  SET title=?, description=?, is_done=?, due_date=?, reminded=?
                WHERE id=? AND user_id=?`,
		input.Title, input.Description, input.IsDone, input.DueDate, input.Reminded, id, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务更新成功"})
}

// DeleteTask DELETE /tasks/{id}
// @Summary      Delete Tasks
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  models.Message
// @Router       /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	_, err = config.DB.Exec("DELETE FROM tasks WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}
