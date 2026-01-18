package handlers

import (
	"github.com/BBloggsbott/task-api/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(taskService *service.TaskService) *gin.Engine {

	router := gin.Default()

	taskHandler := NewTaskHandler(taskService)

	router.GET("/health", taskHandler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.ListTasks)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}

	return router
}
