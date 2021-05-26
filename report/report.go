package report

import (
//	"fmt"
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
	out = config.DB.Where(&model.Task{Com_status: false}).Find(&task)
	c.JSON(http.StatusOK, gin.H{"Remaining Tasks": out.RowsAffected})
}