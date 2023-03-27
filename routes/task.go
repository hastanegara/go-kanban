package routes

import (
	"net/http"
	"project/config"
	"project/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func TaskIndex(c *gin.Context) {
	tasks := []models.Task{}

	config.DB.Preload(clause.Associations).Find(&tasks)

	taskResponse := []models.TaskSimpleResponse{}

	for _, task := range tasks {
		newTask := models.TaskSimpleResponse{
			ID:         task.ID,
			Name:       task.Name,
			Score:      task.Score,
			Duration:   task.Duration,
			StartedAt:  task.StartedAt,
			FinishedAt: task.FinishedAt,
			Status:     task.Status,
		}
		taskResponse = append(taskResponse, newTask)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   taskResponse,
	})
}

func TaskGet(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	data := config.DB.Preload(clause.Associations).First(&task, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data not found",
		})

		c.Abort()

		return
	}

	taskItems := []models.TaskItemResponse{}

	for _, item := range task.TaskItems {
		taskItem := models.TaskItemResponse{
			Name: item.Name,
		}
		taskItems = append(taskItems, taskItem)
	}

	taskUsers := []models.TaskUserResponse{}

	for _, taskUser := range task.TaskUsers {
		var user models.User

		config.DB.First(&user, "id = ?", taskUser.UserID)

		userResponse := models.UserResponse{
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		}
		taskUserResponse := models.TaskUserResponse{
			User:        userResponse,
			Description: taskUser.Description,
			Status:      taskUser.Status,
		}
		taskUsers = append(taskUsers, taskUserResponse)
	}

	taskDocument := models.TaskDocumentResponse{
		Path: task.TaskDocument.Path,
	}
	taskResponse := models.TaskResponse{
		ID:           task.ID,
		Name:         task.Name,
		Description:  task.Description,
		Score:        task.Score,
		Duration:     task.Duration,
		StartedAt:    task.StartedAt,
		FinishedAt:   task.FinishedAt,
		Status:       task.Status,
		TaskItems:    taskItems,
		TaskDocument: taskDocument,
		TaskUsers:    taskUsers,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   taskResponse,
	})
}

func TaskPost(c *gin.Context) {
	var taskRequest models.TaskCreateRequest
	taskItemRequest := []models.TaskItem{}

	c.BindJSON(&taskRequest)

	var taskDocumentRequest models.TaskDocument
	taskDocumentRequest = taskRequest.TaskDocument

	for _, item := range taskRequest.TaskItems {
		taskItemRequest = append(taskItemRequest, item)
	}

	task := models.Task{
		Name:         taskRequest.Name,
		Description:  taskRequest.Description,
		Score:        taskRequest.Score,
		Duration:     taskRequest.Duration,
		TaskItems:    taskItemRequest,
		TaskDocument: taskDocumentRequest,
	}
	createdTask := config.DB.Create(&task)

	if createdTask.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   createdTask.Error.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Created",
	})
}

func TaskPut(c *gin.Context) {
	// 1. Update Existing Task
	id := c.Param("id")
	var task models.Task
	data := config.DB.First(&task, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		c.Abort()

		return
	}

	var taskRequest models.TaskUpdateRequest

	c.BindJSON(&taskRequest)

	taskData := models.Task{
		Name:         taskRequest.Name,
		Description:  taskRequest.Description,
		Score:        taskRequest.Score,
		Duration:     taskRequest.Duration,
		StartedAt:    taskRequest.StartedAt,
		Status:       taskRequest.Status,
		TaskDocument: taskRequest.TaskDocument,
		TaskItems:    taskRequest.TaskItems,
	}
	updatedTask := data.Updates(taskData)

	if updatedTask.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   updatedTask.Error.Error(),
		})

		c.Abort()

		return
	}

	// 2. Wipe (Delete) Related Task Items
	taskItems := []models.TaskItem{}

	config.DB.Where("task_id = ?", id).Delete(&taskItems)

	// 3. Create New Related Task Items
	for _, item := range taskRequest.TaskItems {
		itemData := models.TaskItem{
			TaskID: task.ID,
			Name:   item.Name,
		}
		createTaskItem := config.DB.Create(&itemData)

		if createTaskItem.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal Server Error",
				"error":   createTaskItem.Error.Error(),
			})

			c.Abort()

			return
		}
	}

	// 4. Delete Related Task Document
	taskDocument := models.TaskDocument{}

	config.DB.Where("task_id = ?", id).Delete(&taskDocument)

	// 5. Create New Task Document
	documentData := models.TaskDocument{
		TaskID: task.ID,
		Path:   taskRequest.TaskDocument.Path,
	}
	createTaskDocument := config.DB.Create(&documentData)

	if createTaskDocument.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   createTaskDocument.Error.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated",
	})
}

func TaskDelete(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	data := config.DB.First(&task, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		c.Abort()

		return
	}

	// 1. Delete Related Task Document
	var taskDocument models.TaskDocument
	config.DB.Where("task_id = ?", id).Delete(&taskDocument)

	// 2. Delete Related Task Items
	var taskItem models.TaskItem
	config.DB.Where("task_id = ?", id).Delete(&taskItem)

	// 3. Delete Related Task Users
	var taskUser models.TaskUser
	config.DB.Where("task_id = ?", id).Delete(&taskUser)

	// 4. Delete Task
	config.DB.Delete(&task, id)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Deleted",
	})
}

func TaskStartedAtPatch(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	data := config.DB.First(&task, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		c.Abort()

		return
	}

	var taskRequest models.TaskStartedAtRequest

	c.BindJSON(&taskRequest)

	taskData := models.Task{
		StartedAt: taskRequest.StartedAt,
	}
	updateTask := data.Updates(taskData)

	if updateTask.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   updateTask.Error.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated",
	})
}

func TaskStatustPatch(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	data := config.DB.First(&task, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		c.Abort()

		return
	}

	var taskRequest models.TaskStatusRequest

	c.BindJSON(&taskRequest)

	taskData := models.Task{
		Status: taskRequest.Status,
	}
	updatedTask := data.Updates(taskData)

	if updatedTask.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   updatedTask.Error.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated",
	})
}
