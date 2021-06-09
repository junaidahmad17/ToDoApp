package tasks

import (
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/////////////////////////////////////////////////////////////////////////////////////////////////
func GetUid(c *gin.Context) int {
	y, _ := c.Get("client")
	x, _ := strconv.Atoi(y.(string))
	return x
}

// Listing All Tasks
func GetTasks(c *gin.Context) {
	
	var task []Task
	if e := DB.Where(&Task{Uid: GetUid(c)}).Find(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": e.Error()})
		return 
	}
	c.JSON(http.StatusOK,task)
	
}

// Creating a New Task
func CreateTask(c *gin.Context) {
	DB.AutoMigrate(&Task{})
	
	var task Task
	r,_ := c.Cookie("Token")
	if r == "" {
		c.JSON(http.StatusForbidden,"Forbidden")
		return
	}
	c.BindJSON(&task)

	task.Uid = GetUid(c)
	DB.Create(&task)

	c.JSON(http.StatusOK, "Task Added!")
	
}

// Task Deletion
func DeleteTask(c *gin.Context) {
	
	var task Task
	if e := DB.Where("ID=? AND Uid=?",task.ID,GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	DB.Delete(&task)
	c.JSON(http.StatusOK, "Task Deleted Successfully!")
}

// Editing a Tasking Using ID as Key 
func EditTask(c *gin.Context) {

	var task Task
	if e := DB.Where("ID=? AND Uid=?",task.ID,GetUid(c)).First(&task).Error; e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Task does not exist!"})
		return 
	}

	
	var input UpdateTask
	c.BindJSON(&input)
	DB.Model(&task).Updates(Task{Title: input.Title, Description: input.Description, Com_status: input.Com_status})
	c.JSON(http.StatusOK, "Task Modified Successfully!")
}

// Clearing Whole DB
func DeleteAll(c *gin.Context) {
	var task0 []Task
	DB.Find(&task0)
	DB.Delete(&task0)
	c.JSON(http.StatusOK, "All Entries Deleted!")
}

////////////////////////////           Reports                //////////////////////////////////////////////
// Listing Basic Stats
func CountTask(c *gin.Context) {
	var task []Task
	// Total Tasks 
	out := DB.Where("Uid=?",GetUid(c)).Find(&task)
	c.JSON(http.StatusOK, gin.H{"Total Tasks": out.RowsAffected})
	// Completed Tasks 
	out = DB.Where(&Task{Com_status: true,Uid: GetUid(c) }).Find(&task)
	c.JSON(http.StatusOK, gin.H{"Completed Tasks": out.RowsAffected})
	// Remaining Tasks
	out = DB.Where(&Task{},"Com_status").Find(&task)
	c.JSON(http.StatusOK, gin.H{"Remaining Tasks": out.RowsAffected})
}

// Incomplete Tasks Past Due Dates
func MissedTasks(c *gin.Context) {

	t := time.Now()
	count := 0

	rows, _ := DB.Model(&Task{}).Where("Com_status = ? AND Uid =?", false,GetUid(c)).Rows()
	var task Task
	for rows.Next() {
		DB.ScanRows(rows,&task)
		
		if task.Due_DT.Sub(t) < 0 {
			count = count + 1
		}
	}

	c.JSON(http.StatusOK, gin.H{"Missed Deadlines":count})
}