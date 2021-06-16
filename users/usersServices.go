package users

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todoapp/auth"
	"todoapp/email"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func encodeURL(email string) string {
	encodedURL := base64.URLEncoding.EncodeToString([]byte(email))
	return encodedURL
}

func VerifyEmail(x string, c *gin.Context) bool {
	var user User
	out := UDB.Where("Email=?", x).First(&user)
	flag := out.RowsAffected==0
	return flag
}

func CreateUser(c *gin.Context) {
	UDB.AutoMigrate(&User{})
	var user User
	c.BindJSON(&user)
	
	ePassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	user.Password = string(ePassword)
	

	if VerifyEmail(user.Email, c) {
		
		UDB.Create(&user)
		subject := "ToDo Account Resgisteration"
		content := "Dear "+user.Username+",\n\n"
		content = content+"A ToDo account has been registered using your email! Please click on the link below to confirm.\n"
		content = content+"/verify/"+encodeURL(user.Email)
		content = content+"\n\nRegards,\nToDo Team"
		
		email.SendEmail(user.Email,subject ,content)
		c.JSON(http.StatusCreated ,"Account Registered")
		return 
	}
	c.JSON(http.StatusNotAcceptable, gin.H{"Error":"Email already registered"} )
}

func VerifyUser(c *gin.Context) {
	decodedToken, _ :=  base64.URLEncoding.DecodeString(c.Param("token"))
	var user *User
	out := UDB.Where(&User{Email:string(decodedToken)}).First(&user)
	if out.RowsAffected ==0 {
		c.JSON(http.StatusNotFound, "Invalid verification link!")
		return
	}
	user.EmailVerified = true 
	UDB.Save(&user)	
	c.JSON(http.StatusOK, "Email Verified Successfully!")
}

func Login(c *gin.Context) {
	
	var input User
	var user User
	c.BindJSON(&input)

	out := UDB.Where("Email = ?", input.Email).Find(&user)
	if out.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, "1Invalid Email or Password!")
		return
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusUnauthorized, "2Invalid Email or Password!")
		return
	}
	if !user.EmailVerified {
		c.JSON(http.StatusExpectationFailed, "Please verify your email address first!")
		return
	}

	token,err := auth.CreateJWT(strconv.FormatUint(uint64(user.ID), 10), user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("Token",token,1000000,"/","",false,false)

	c.JSON(http.StatusOK, "User logged in successfully")

}

func Logout(c *gin.Context) {
	r,_ := c.Cookie("Token")
	if r == "" {
		c.JSON(http.StatusForbidden,"Forbidden")
		return
	}
	
	c.SetCookie("Token","",-1,"/","",false,false)

	c.JSON(http.StatusOK,"Logged out successfully!")
	
}

func ForgotPass(c *gin.Context) {
	var input User
	c.BindJSON(&input)
	
	if !VerifyEmail(input.Email, c) {
		s := input.Email+","
		s = s + time.Now().String()

		subject := "Password Reset Request"
		content := "Hi!\n" 
		content = content+"Please follow the reset link given below to reset your password!\n\n"
		content = content+"/resetlink/"+encodeURL(s)
		content = content+"\n\nRegards,\nToDo Team"
		
		email.SendEmail(input.Email,subject ,content)
		c.JSON(http.StatusCreated ,"A reset link has been sent to your email. Please check your inbox.")
		return
	}
	
	c.JSON(http.StatusNotFound, "Unregistered Email!")
}

func ResetLink(c *gin.Context) {
	decodedToken, _ :=  base64.URLEncoding.DecodeString(c.Param("token"))
	decodedTokens := string(decodedToken)
	email := strings.Split(decodedTokens, ",")[0]

	var user User

	out := UDB.Where(&User{Email:email}).First(&user)
	if out.RowsAffected==0 {
		c.JSON(http.StatusNotFound, "Invalid reset link")
		return
	}
	var input User
	err := c.BindJSON(&input)
	if err!= nil {
		fmt.Println("Error: ", err.Error())	
	}
	// Hashing Password
	ePassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	user.Password = string(ePassword)
	
	UDB.Save(&user)
	c.JSON(http.StatusOK, "Password changed successfully!")
}