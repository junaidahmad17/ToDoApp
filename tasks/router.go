package tasks

import (
	"github.com/gin-gonic/gin"
	"todoapp/auth"
)


func InitTasks(r *gin.Engine) {
	
	// tasks

	//r.GET("/tasks", auth(), GetTasks)
	r.GET("/tasks", auth.IsAuthorized(), GetTasks)
	 
	r.POST("/tasks", auth.IsAuthorized(), CreateTask)

	r.PATCH("/tasks/:id", auth.IsAuthorized(), EditTask)

	r.DELETE("/tasks/:id", auth.IsAuthorized(), DeleteTask)

	r.DELETE("/all", auth.IsAuthorized(), DeleteAll)

	// Reports

	r.GET("/report/", auth.IsAuthorized(), CountTask)

	r.GET("/report/mt", auth.IsAuthorized(), MissedTasks)

}