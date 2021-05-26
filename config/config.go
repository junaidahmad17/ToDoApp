package config 

import (
	"gorm.io/gorm"
  	"gorm.io/driver/sqlite"
	"fmt"
	"os"
)
var DB *gorm.DB
var err error

func SetupDB() {
	if DB == nil{
		e := os.Remove("C:\\Users\\Junaid Ahmad (WORK)\\Desktop\\GO\\todoapp\\ToDo.db")
		if e != nil{
			fmt.Println("Error:  ", e)
		}

		DB, err = gorm.Open(sqlite.Open("ToDo.db"), &gorm.Config{})
		if err != nil {
		fmt.Println("Status:", "Error_Get!!!")
		}
	}
}