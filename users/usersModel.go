package users

import(
	"gorm.io/gorm"
  	"gorm.io/driver/sqlite"
	"fmt"
	"os"
)

type User struct {

	ID int
	Username string 
	Password string
	Email string
	//TaskU []tasks.Task //`gorm:"ForeignKey:UserID"`
}

var UDB *gorm.DB
var Uerr error

func SetupUDB() {
	flag := true
	
	if UDB == nil {
	
		if flag == false {
			e := os.Remove("C:\\Users\\Junaid Ahmad (WORK)\\Desktop\\GO\\newtodo\\todoapp\\Users.db")
			fmt.Println("Erased------------------------------")
			if e != nil {
				fmt.Println("Error:  ", e)
			}
		}
		
		UDB, Uerr = gorm.Open(sqlite.Open("Users.db"), &gorm.Config{})
		
		if Uerr != nil {
		fmt.Println("Status:", "Error_Get!!!")
		}
	}
}
