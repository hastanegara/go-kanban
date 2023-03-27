package routes

import (
	"net/http"
	"project/config"
	"project/models"
	"time"

	"github.com/gin-gonic/gin"
)

func TaskUserPost(c *gin.Context) {
	// 1. Create Task User
	var taskUserRequest models.TaskUserRequest

	c.BindJSON(&taskUserRequest)

	taskUser := models.TaskUser{
		TaskID:      taskUserRequest.TaskID,
		UserID:      taskUserRequest.UserID,
		Description: taskUserRequest.Description,
		Status:      0,
	}
	createdTaskUser := config.DB.Create(&taskUser)

	if createdTaskUser.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   createdTaskUser.Error.Error(),
		})

		c.Abort()

		return
	}

	// 2. Update Task Started at and Status = 1
	var task models.Task
	taskData := models.Task{
		StartedAt: time.Now().String(),
		Status:    1,
	}
	updatedTask := config.DB.Find(&task, "id = ?", taskUserRequest.TaskID)

	updatedTask.Updates(taskData)

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "New Task User created",
	})
}

func TaskUserPatch(c *gin.Context) {
	// 1. Search existing Task User
	id := c.Param("id")
	var taskUser models.TaskUser
	data := config.DB.First(&taskUser, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		c.Abort()

		return
	}

	if taskUser.Status == 1 {
		c.JSON(http.StatusConflict, gin.H{
			"status":  http.StatusConflict,
			"message": "Conflict. Data already assigned",
		})

		c.Abort()

		return
	}

	// 2. Update Task User Status = 1
	taskUserData := models.TaskUser{
		Status: 1,
	}
	updatedTaskUser := data.Updates(taskUserData)

	if updatedTaskUser.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   updatedTaskUser.Error.Error(),
		})

		c.Abort()

		return
	}

	// 3. Update Task Finished at and Status = 2
	var task models.Task
	taskData := config.DB.First(&task, "id = ?", taskUser.TaskID)
	updatedTask := models.Task{
		Status:     2,
		FinishedAt: time.Now().String(),
	}

	taskData.Updates(updatedTask)

	// 4. Create New User Score or Update Exsisting User Score
	var userScore models.UserScore
	scoreData := config.DB.First(&userScore, "user_id = ?", taskUser.UserID)

	if scoreData.RowsAffected == 0 {
		createdScore := models.UserScore{
			UserID: taskUser.UserID,
			Score:  task.Score,
		}

		config.DB.Create(&createdScore)
	} else {
		score := userScore.Score + task.Score
		updatedScore := models.UserScore{
			Score: score,
		}

		scoreData.Updates(updatedScore)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Status Task Updated",
	})
}
