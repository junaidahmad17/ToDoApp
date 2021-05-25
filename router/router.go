package router // Controller 

import (
	"todoapp/handler"
	"github.com/gin-gonic/gin"
)
   

func SetupRouter() *gin.Engine {
	r := gin.Default()
	
	r.GET("/tasks", handler.GetTasks)
	 
	r.POST("/tasks", handler.CreateTask)

	r.PATCH("/tasks/:id", handler.EditTask)
	
	r.DELETE("/tasks/:id", handler.DeleteTask)

	r.Run(":8080") // Main File 

	return r
}