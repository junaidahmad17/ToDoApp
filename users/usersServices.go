package users

import (
	"fmt"
	"net/http"
	//"fmt"
	"todoapp/auth"
	"strconv"
	"github.com/gin-gonic/gin"
)
func VerifyEmail(x string) bool {
	var user User
	if e := UDB.Where("Email=?", x).First(&user).Error; e == nil {
		return false
	}
	return true
}

func CreateUser(c *gin.Context) {
	
	UDB.AutoMigrate(&User{})
	var err error
	var user User
	c.BindJSON(&user)
	

	if VerifyEmail(user.Email) {
		if err = UDB.Create(user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error":"Failed to open db"} )
					
		} else {
			c.JSON(http.StatusCreated ,"Account Registered")
		}
}
}


func Login(c *gin.Context) {
	//fmt.Println("1--------------------------------------------------------------------")
	var reqBody User
	var user User
	c.BindJSON(&reqBody)
	//fmt.Println("2--------------------------------------------------------------------")
	//err := json.NewDecoder(r.Body).Decode(&reqBody)	fmt.Println("3--------------------------------------------------------------------")
	
	//fmt.Println("4--------------------------------------------------------------------")
	//c.JSON(http.StatusOK, "Email: "+reqBody.Email+", Password: "+reqBody.Password)
	//fmt.Println("5--------------------------------------------------------------------")
	
	out := UDB.Where("Email = ? AND Password = ?", reqBody.Email, reqBody.Password).Find(&user)
	if out.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, "Invalid Email or Password!")
		return
	}
	
	token, err := auth.CreateJWT(strconv.Itoa(user.ID), user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(token)
	c.SetCookie("Token",token,1000000,"/","",false,false)

	c.JSON(http.StatusOK, "User logged in successfully")

}
func Msg(c *gin.Context) {
	c.JSON(200,"Secret Code")
}
func Logout(c *gin.Context) {
	r,_ := c.Cookie("Token")
	if r == "" {
		c.JSON(http.StatusForbidden,"Forbidden")
		return
	}
	
	c.SetCookie("Token","",-1,"/","",false,false)

	c.JSON(http.StatusOK,"Logged out successfully!")
	//c.JSON(200,r)
}

/*func logout(w http.ResponseWriter, r *http.Request) {
	
	_, err := auth.GetJWT(r)
	if err != nil {
		utils.JSONMsg(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	accessCookie := http.Cookie{Name: "accessCookie", Value: "", MaxAge: -1}
	http.SetCookie(w, &accessCookie)
	utils.JSONMsg(w, "User logged out successfully", http.StatusOK)

}*/
