package users

import (
	"github.com/gin-gonic/gin"
	"todoapp/auth"
)


func InitUsers(r *gin.Engine) {

	r.POST("/register", CreateUser)
	r.POST("/login", Login)
	r.GET("/logout", auth.IsAuthorized(),Logout)
}
