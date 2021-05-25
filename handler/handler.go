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

type UpdateTask struct {
	Title string       `json:"title"`
	Description string `json:"description"`
  	Com_status bool	`json:"com_status"`
}

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
	fmt.Println(task0.ID)
	c.JSON(http.StatusOK, "Task Added!")
	
}

func DeleteTask(c *gin.Context) {
	
	setupDB()

	var task0 model.Task
	if e := db.Where("id=?",c.Param("id")).First(&task0).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	db.Delete(&task0)
	c.JSON(http.StatusOK, "Task Deleted Successfully!")
}

func EditTask(c *gin.Context) {
	
	setupDB()

	var task0 model.Task
	if e := db.Where("id=?",c.Param("id")).First(&task0).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	
	var input UpdateTask
	c.BindJSON(&input)
	fmt.Println("Com_status", input.Com_status)
	db.Model(&task0).Updates(model.Task{Title: input.Title, Description: input.Description, Com_status: input.Com_status})
	c.JSON(http.StatusOK, "Task Modified Successfully!")
}