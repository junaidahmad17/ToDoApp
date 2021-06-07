package users

import (
	//"todoapp/auth"

	"github.com/gin-gonic/gin"
	"todoapp/auth"
)


func InitUsers(r *gin.Engine) {

	r.POST("/register", CreateUser)
	r.POST("/login", Login)
	r.GET("/logout", auth.IsAuthorized(),Logout)
	r.GET("/code", auth.IsAuthorized(), Msg)
	//http.Handle("/logout", aServices.IsAuthorized(http.HandlerFunc(logout)))
}
