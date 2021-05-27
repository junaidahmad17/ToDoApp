package handler

import (
	"fmt"
	"todoapp/model"
	"github.com/gin-gonic/gin"
	"todoapp/config"
	"net/http"
)

type UpdateTask struct {
	Title string       `json:"title"`
	Description string `json:"description"`
  	Com_status bool	`json:"com_status"`
}


func GetTasks(c *gin.Context) {
	
	var task0 []model.Task
	config.DB.Find(&task0)
	c.JSON(http.StatusOK,task0)
	
}

func CreateTask(c *gin.Context) {

	config.DB.AutoMigrate(&model.Task{})

	var task0 model.Task
	c.BindJSON(&task0)
	config.DB.Create(&task0)
	fmt.Println(task0.ID)
	c.JSON(http.StatusOK, "Task Added!")
	
}

func DeleteTask(c *gin.Context) {
	
	var task0 model.Task
	if e := config.DB.Where("id=?",c.Param("id")).First(&task0).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	config.DB.Delete(&task0)
	c.JSON(http.StatusOK, "Task Deleted Successfully!")
}

func EditTask(c *gin.Context) {

	var task0 model.Task
	if e := config.DB.Where("id=?",c.Param("id")).First(&task0).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	
	var input UpdateTask
	c.BindJSON(&input)
	fmt.Println("Com_status", input.Com_status)
	config.DB.Model(&task0).Updates(model.Task{Title: input.Title, Description: input.Description, Com_status: input.Com_status})
	c.JSON(http.StatusOK, "Task Modified Successfully!")
}

func DeleteAll(c *gin.Context) {
	var task0 []model.Task
	config.DB.Find(&task0)
	config.DB.Delete(&task0)
	c.JSON(http.StatusOK, "All Entries Deleted!")
}