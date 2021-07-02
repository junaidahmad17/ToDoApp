package users

import (
	"github.com/gin-gonic/gin"
	"todoapp/auth"
)


func InitUsers(r *gin.Engine) {

	r.POST("/register", CreateUser)
	r.POST("/login", Login)
	r.GET("/logout", auth.IsAuthorized(),Logout)

	r.GET("/verify/:token", VerifyUser)

	r.POST("/reset", ForgotPass)
	r.POST("/resetlink/:token", ResetLink)


	///  OAuth2
	
	r.GET("/gmail", GmailLogin)
	r.GET("/auth/google", Redirect)
	r.GET("/auth/google/callback", Callback)
}
