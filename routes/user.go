package routes

import (
	"net/http"
	"project/auth"
	"project/config"
	"project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func UserIndex(c *gin.Context) {
	users := []models.User{}

	config.DB.Find(&users)

	userResponse := []models.UserResponse{}

	for _, user := range users {
		newUserScore := models.UserScoreResponse{
			Score: user.UserScore.Score,
		}
		newUser := models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      user.Role,
			UserScore: newUserScore,
		}

		userResponse = append(userResponse, newUser)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   userResponse,
	})
}

func UserGet(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	data := config.DB.Preload(clause.Associations).First(&user, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		c.Abort()

		return
	}

	userScoreResponse := models.UserScoreResponse{
		Score: user.UserScore.Score,
	}
	var userTasks []models.UserTaskResponse

	for _, task := range user.UserTasks {
		var newTask models.Task

		config.DB.First(&newTask, "id = ?", task.TaskID)

		taskResponse := models.TaskResponse{
			Name: newTask.Name,
		}
		newUserTask := models.UserTaskResponse{
			Task:        taskResponse,
			Description: task.Description,
			Status:      task.Status,
		}
		userTasks = append(userTasks, newUserTask)
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Role:      user.Role,
		UserScore: userScoreResponse,
		UserTasks: userTasks,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   userResponse,
	})
}

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad Request",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	err := user.HashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	createUser := config.DB.Create(&user)

	if createUser.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   createUser.Error.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Created",
	})
}

func GenerateToken(c *gin.Context) {
	request := models.TokenRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Bad Request",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	user := models.User{}
	checkUsername := config.DB.Where("username = ?", request.Username).First(&user)

	if checkUsername.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Username not found",
			"error":   checkUsername.Error.Error(),
		})

		c.Abort()

		return
	}

	credentialError := user.CheckPassword(request.Password)

	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Password not match",
			"error":   credentialError.Error(),
		})

		c.Abort()

		return
	}

	tokenString, err := auth.GenerateToken(user.Username, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  tokenString,
	})
}

func UserPut(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	data := config.DB.First(&user, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		c.Abort()

		return
	}

	var userRequest models.UserRequest

	c.BindJSON(&userRequest)

	bytes, _ := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 12)
	password := string(bytes)
	userData := models.User{
		Username: userRequest.Username,
		Password: password,
		FullName: userRequest.FullName,
	}
	updatedUser := data.Updates(userData)

	if updatedUser.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal Server Error",
			"error":   updatedUser.Error.Error(),
		})

		c.Abort()

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Updated",
	})
}

func UserDelete(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	data := config.DB.First(&user, "id = ?", id)

	if data.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Data Not Found",
		})

		return
	}

	// 1. Delete Related User Score
	var userScore models.UserScore
	config.DB.Where("user_id = ?", id).Delete(&userScore)

	// 2. Delete Related Task Users
	var taskUser models.TaskUser
	config.DB.Where("user_id = ?", id).Delete(&taskUser)

	// 3. Delete User
	config.DB.Delete(&user, id)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Deleted",
	})
}
