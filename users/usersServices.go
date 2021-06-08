package users

import (
	"net/http"
	"strconv"
	"todoapp/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)
func VerifyEmail(x string, c *gin.Context) bool {
	var user User
	out := UDB.Where("Email=?", x).First(&user)
	flag := out.RowsAffected==0
	if !flag {
		c.JSON(http.StatusNotAcceptable, gin.H{"Error":"Email already registered"} )
	}
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
		Count = Count + 1
		user.IDU = Count
		UDB.Create(user)
		c.JSON(http.StatusCreated ,"Account Registered")
	}
}


func Login(c *gin.Context) {
	
	var input User
	var user User
	c.BindJSON(&input)

	out := UDB.Where("Email = ?", input.Email).Find(&user)
	if out.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, "Invalid Email or Password!")
		return
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusUnauthorized, "Invalid Email or Password!")
		return
	}

	token,err := auth.CreateJWT(strconv.FormatUint(uint64(user.IDU), 10), user.Password)
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
