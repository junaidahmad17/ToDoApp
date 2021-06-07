package tasks

import (
	"fmt"
	//"todoapp/users"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////////////////////

// Listing All Tasks
func GetTasks(c *gin.Context) {
	
	var task0 []Task
	DB.Find(&task0)
	//fmt.Println(s)
	c.JSON(http.StatusOK,task0)
	
}

// Creating a New Task
func CreateTask(c *gin.Context) {
	
	DB.AutoMigrate(&Task{})
	//c.Set("UserID", Uid) 
	//x = c.Get("UserID")
	var task0 Task
	r,_ := c.Cookie("Token")
	if r == "" {
		c.JSON(http.StatusForbidden,"Forbidden")
		return
	}
	c.BindJSON(&task0)
	y, _ := c.Get("client")
	fmt.Println(y)
	x, _ := strconv.Atoi(y.(string))
	task0.Uid = x 
	DB.Create(&task0)
	fmt.Println(task0.ID)
	c.JSON(http.StatusOK, "Task Added!")
	
}

// Task Deletion
func DeleteTask(c *gin.Context) {
	
	var task0 Task
	if e := DB.Where("id=?",c.Param("id")).First(&task0).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	DB.Delete(&task0)
	c.JSON(http.StatusOK, "Task Deleted Successfully!")
}

// Editing a Tasking Using ID as Key 
func EditTask(c *gin.Context) {

	var task0 Task
	if e := DB.Where("id=?",c.Param("id")).First(&task0).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	
	var input UpdateTask
	c.BindJSON(&input)
	fmt.Println("Com_status", input.Com_status)
	DB.Model(&task0).Updates(Task{Title: input.Title, Description: input.Description, Com_status: input.Com_status})
	c.JSON(http.StatusOK, "Task Modified Successfully!")
}

// Clearing Whole DB
func DeleteAll(c *gin.Context) {
	var task0 []Task
	DB.Find(&task0)
	DB.Delete(&task0)
	c.JSON(http.StatusOK, "All Entries Deleted!")
}