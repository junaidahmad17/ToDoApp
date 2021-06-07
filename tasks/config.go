package tasks

import (
	"gorm.io/gorm"
  	"gorm.io/driver/sqlite"
	//"todoapp/userModel"
	"fmt"
	"os"
)
var DB *gorm.DB
var err error

func SetupDB() {
	flag := false
	
	if DB == nil {
	
		if flag == false {
			e := os.Remove("C:\\Users\\Junaid Ahmad (WORK)\\Desktop\\GO\\newtodo\\todoapp\\ToDo.db")
			if e != nil {
				fmt.Println("Error:  ", e)
			}
		}
		DB, err = gorm.Open(sqlite.Open("ToDo.db"), &gorm.Config{})
		if err != nil {
		fmt.Println("Status:", "Error_Get!!!")
		}
	}
}

