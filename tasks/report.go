package tasks

import (
	//"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Listing Basic Stats
func CountTask(c *gin.Context) {
	var task []Task
	// Total Tasks 
	out := DB.Find(&task)
	c.JSON(http.StatusOK, gin.H{"Total Tasks": out.RowsAffected})
	// Completed Tasks 
	out = DB.Where(&Task{Com_status: true}).Find(&task)
	c.JSON(http.StatusOK, gin.H{"Completed Tasks": out.RowsAffected})
	// Remaining Tasks
	out = DB.Where(&Task{},"Com_status").Find(&task)
	c.JSON(http.StatusOK, gin.H{"Remaining Tasks": out.RowsAffected})
}

// Incomplete Tasks Past Due Dates
func MissedTasks(c *gin.Context) {

	t := time.Now()
	count := 0

	rows, _ := DB.Model(&Task{}).Where("Com_status = ?", false).Rows()
	var task0 Task
	for rows.Next() {
		DB.ScanRows(rows,&task0)
		
		if task0.Due_DT.Sub(t) < 0 {
			count = count + 1
		}
	}

	c.JSON(http.StatusOK, gin.H{"Missed Deadlines":count})
}