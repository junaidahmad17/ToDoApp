package users

import(
	"gorm.io/gorm"
  	"gorm.io/driver/sqlite"
	"fmt"
	"os"
)
var Count uint
type User struct {

	ID uint `json:"id"`
	Username string 
	Password string
	Email string
	EmailVerified bool
}

var UDB *gorm.DB
var Uerr error
//Allow insecure apps
func SetupUDB() {
	flag := true
	Count = 0	
	if UDB == nil {
		println("Value: ", os.Getenv("DBADD"))
		if !flag {
			e := os.Remove(os.Getenv("DBADD")+"Users.db")
			if e != nil {
				fmt.Println("Error:  ", e.Error())
			}
		}
		
		UDB, Uerr = gorm.Open(sqlite.Open("Users.db"), &gorm.Config{})
		
		if Uerr != nil {
		fmt.Println("Status:", "Failed to open database!")
		}
	}
}
