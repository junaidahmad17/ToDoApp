package users

import(
	"gorm.io/gorm"
  	"gorm.io/driver/sqlite"
	"fmt"
	"os"
)
var Count uint
type User struct {

	IDU uint
	Username string 
	Password string
	Email string
}

var UDB *gorm.DB
var Uerr error

func SetupUDB() {
	flag := true
	Count = 0	
	if UDB == nil {
	
		if !flag {
			e := os.Remove("C:\\Users\\Junaid Ahmad (WORK)\\Desktop\\GO\\newtodo\\todoapp\\Users.db")
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
