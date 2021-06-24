package users

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
				log.Fatalln(e)
			}
		}
		
		UDB, Uerr = gorm.Open(sqlite.Open("Users.db"), &gorm.Config{})
		
		if Uerr != nil {
		log.Fatalln(Uerr)
		}
	}
}
