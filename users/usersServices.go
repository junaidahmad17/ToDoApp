package users

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
	"log"
	"os"
	"math/rand"
	"todoapp/auth"
	"todoapp/email"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	//   OAuth2
	"html/template"
	"github.com/markbates/goth"
    "github.com/markbates/goth/gothic"
    "github.com/markbates/goth/providers/google"
    "github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	key := os.Getenv("MYCODE")  
  	maxAge := 86400 * 30  
  	isProd := false       

  	store := sessions.NewCookieStore([]byte(key))
  	store.MaxAge(maxAge)
  	store.Options.Path = "/"
  	store.Options.HttpOnly = true 
  	store.Options.Secure = isProd

  	gothic.Store = store

  	goth.UseProviders(
    	google.New(os.Getenv("Client_ID"), os.Getenv("Client_Secret"), "http://localhost"+os.Getenv("PORT")+"/auth/google/callback", "email", "profile"),
  	)
	  
}
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
func allowCreateUser(user User,c *gin.Context) {
	UDB.Create(&user)
		subject := "ToDo Account Resgisteration"
		content := "Dear "+user.Username+",\n\n"
		content = content+"A ToDo account has been registered using your email! Please click on the link below to confirm.\n"
		content = content+"/verify/"+encodeURL(user.Email)
		content = content+"\n\nRegards,\nToDo Team"
		
		email.SendEmail(user.Email,subject ,content)
		c.JSON(http.StatusCreated ,"Account Registered")
}
func CreateUser(c *gin.Context) {
	UDB.AutoMigrate(&User{})
	var user User
	c.BindJSON(&user)
	
	ePassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
	}
	user.Password = string(ePassword)
	

	if VerifyEmail(user.Email, c) {
		
		allowCreateUser(user, c)
		return 
	}
	c.JSON(http.StatusNotAcceptable, gin.H{"error":"email already registered"} )
}

func VerifyUser(c *gin.Context) {
	decodedToken, _ :=  base64.URLEncoding.DecodeString(c.Param("token"))
	var user *User
	out := UDB.Where(&User{Email:string(decodedToken)}).First(&user)
	if out.RowsAffected ==0 {
		c.JSON(http.StatusNotFound, gin.H{"error":"invalid verification link!"})
		return
	}
	user.EmailVerified = true 
	UDB.Save(&user)	
	c.JSON(http.StatusOK, gin.H{"msg":"email verified successfully!"})
}

func allowLogin(x uint, y string,c *gin.Context) {
	token,err := auth.CreateJWT(strconv.FormatUint(uint64(x), 10), y)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("Token",token,1000000,"/","",false,false)

	c.JSON(http.StatusOK, gin.H{"msg":"user logged in successfully"})

}
func Login(c *gin.Context) {
	
	var input User
	var user User
	c.BindJSON(&input)

	out := UDB.Where("Email = ?", input.Email).Find(&user)
	if out.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error":"invalid email or password!"})
		return
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"ivalid email or password!"})
		return
	}
	if !user.EmailVerified {
		c.JSON(http.StatusExpectationFailed, gin.H{"error":"please verify your email address first!"})
		return
	}
	allowLogin(user.ID,user.Password,c)
	
}

func Logout(c *gin.Context) {
	
	r,_ := c.Cookie("Token")
	if r == "" {
		c.JSON(http.StatusForbidden,gin.H{"error":"forbidden"})
		return
	}
	
	c.SetCookie("Token","",-1,"/","",false,false)

	c.JSON(http.StatusOK,gin.H{"msg":"logged out successfully!"})
	
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
	
	c.JSON(http.StatusNotFound, gin.H{"msg": "unregistered email!"})
}

func ResetLink(c *gin.Context) {
	decodedToken, _ :=  base64.URLEncoding.DecodeString(c.Param("token"))
	decodedTokens := string(decodedToken)
	email := strings.Split(decodedTokens, ",")[0]

	var user User

	out := UDB.Where(&User{Email:email}).First(&user)
	if out.RowsAffected==0 {
		c.JSON(http.StatusNotFound, gin.H{"msg":"invalid reset link"})
		return
	}
	var input User
	err := c.BindJSON(&input)
	if err!= nil {
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})	
	}
	// Hashing Password
	ePassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	user.Password = string(ePassword)
	
	UDB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"error":"password changed successfully!"})
}


////////////////////////////////////////        OAuth2        //////////////////////////////////////////////

func Redirect(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func Callback(c *gin.Context)  {
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
    user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
      log.Println(gin.H{"error": err})
      return
    }

	t, _ := template.ParseFiles("templates/success.html")

	if !VerifyEmail(user.Email, c) {
		var check User
		UDB.Where("Email = ?", user.Email).Find(&check)
		allowLogin(check.ID,check.Password,c)
	} else {
		var u User
		u.Email = user.Email
		u.Username = user.Name
		u.Password = strconv.Itoa(rand.Int())
		u.EmailVerified = true
		allowCreateUser(u, c)
	}
	t.Execute(c.Writer, user)
}

func GmailLogin(c *gin.Context) {
	t, _ := template.ParseFiles("templates/index.html")
    t.Execute(c.Writer, false)
}
