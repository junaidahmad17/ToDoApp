package tasks

import (
	"github.com/gin-gonic/gin"
	"todoapp/auth"
)


func InitTasks(r *gin.Engine) {
	
	r.GET("/tasks", auth.IsAuthorized(), GetTasks)
	 
	r.POST("/tasks", auth.IsAuthorized(), CreateTask)

	r.PATCH("/tasks/:id", auth.IsAuthorized(), EditTask)

	r.DELETE("/tasks/:id", auth.IsAuthorized(), DeleteTask)

	r.DELETE("/all", auth.IsAuthorized(), DeleteAll)

	// Reports

	r.GET("/report/", auth.IsAuthorized(), CountTask)

	r.GET("/report/mt", auth.IsAuthorized(), MissedTasks)

	// Attachment 

	r.PUT("/uploadfile/:id",auth.IsAuthorized(), attachFile)
	r.DELETE("/delfile/:id",auth.IsAuthorized(), deleteFile)
	r.GET("/download/:id",auth.IsAuthorized(), downloadFile)

}