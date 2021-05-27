package report

import (
	//"fmt"
	"time"
	"todoapp/model"
	"github.com/gin-gonic/gin"
	"todoapp/config"
	"net/http"
)

func CountTask(c *gin.Context) {
	var task []model.Task
	out := config.DB.Find(&task)
	c.JSON(http.StatusOK, gin.H{"Total Tasks": out.RowsAffected})
	out = config.DB.Where(&model.Task{Com_status: true}).Find(&task)
	c.JSON(http.StatusOK, gin.H{"Completed Tasks": out.RowsAffected})
	out = config.DB.Where(&model.Task{},"Com_status").Find(&task)
	c.JSON(http.StatusOK, gin.H{"Remaining Tasks": out.RowsAffected})
}

func MissedTasks(c *gin.Context) {

	t := time.Now()
	count := 0

	rows, _ := config.DB.Model(&model.Task{}).Where("Com_status = ?", false).Rows()
	var task0 model.Task
	for rows.Next() {
		config.DB.ScanRows(rows,&task0)
		
		if task0.Due_DT.Sub(t) < 0 {
			count = count + 1
		}
	}

	c.JSON(http.StatusOK, gin.H{"Missed Deadlines":count})
}