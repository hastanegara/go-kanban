package main

import (
	"net/http"
	"project/config"
	"project/middlewares"
	"project/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DBConnect()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", GetHome)

		user := v1.Group("/user")
		{
			user.POST("/register", routes.Register)
			user.POST("/login", routes.GenerateToken)

			user.Use(middlewares.IsLeader())
			{
				user.GET("/", routes.UserIndex)
				user.GET("/:id", routes.UserGet)
				user.PUT("/:id", routes.UserPut)
				user.DELETE("/:id", routes.UserDelete)
			}
		}

		task := v1.Group("/task")
		{
			task.Use(middlewares.Auth())
			{
				task.GET("/", routes.TaskIndex)
				task.GET("/:id", routes.TaskGet)
				task.PATCH("/started_at/:id", routes.TaskStartedAtPatch)
				task.PATCH("/status/:id", routes.TaskStatustPatch)
			}

			task.Use(middlewares.IsLeader())
			{
				task.POST("/", routes.TaskPost)
				task.PUT("/:id", routes.TaskPut)
				task.DELETE("/:id", routes.TaskDelete)
			}

		}

		taskUser := v1.Group("/task-user").Use(middlewares.Auth())
		{
			taskUser.POST("/", routes.TaskUserPost)
			taskUser.PATCH("/:id", routes.TaskUserPatch)
		}
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func GetHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Message": "Welcome to Kanban Board Project developed by Heidar Hastanegara",
	})
}

// https://github.com/hastanegara/go-kanban
// https://api.postman.com/collections/1421695-bc73f7ba-baa8-427c-8878-5f2fbc342b8f?access_key=PMAT-01GWHENMRTZ1PHDHWFEF7SEG3R
