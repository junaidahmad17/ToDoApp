package handler

import (
	"fmt"
	"todoapp/model"
	"gorm.io/gorm"
  	"gorm.io/driver/sqlite"
	"github.com/gin-gonic/gin"
	"os"
	"net/http"
)
var db *gorm.DB
var err error

func setupDB(){
	if db == nil{
		e := os.Remove("C:\\Users\\Junaid Ahmad (WORK)\\Desktop\\GO\\todoapp\\ToDo.db")
		if e!=nil{
			fmt.Println("Error:  ", e)
		}

		db, err = gorm.Open(sqlite.Open("ToDo.db"), &gorm.Config{})
		if err != nil {
		fmt.Println("Status:", "Error_Get!!!")
		}
	}
}
func GetTasks(c *gin.Context) {
	setupDB()

	var task0 []model.Task
	db.Find(&task0)
	c.JSON(http.StatusOK,task0)
	
}

func CreateTask(c *gin.Context) {
	setupDB()

	db.AutoMigrate(&model.Task{})

	var task0 model.Task
	c.BindJSON(&task0)
	db.Create(&task0)
	c.JSON(http.StatusOK, "Task Added!")
	
}
