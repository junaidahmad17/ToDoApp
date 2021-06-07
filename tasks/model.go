package tasks

import (
	"time"
)
type Task struct {

	// user id foreign key int 
	ID uint            `json:"id"`
	Title string       `json:"title"`
	Description string `json:"description"`
	Create_DT  time.Time `json:"create_DT"`
	Due_DT time.Time	`json:"due_DT"`
  	Com_status bool	`json:"com_status"`
  	Com_DT time.Time	`json:"com_DT"`
	Uid int
}

type UpdateTask struct {
	Title string       `json:"title"`
	Description string `json:"description"`
  	Com_status bool	`json:"com_status"`
}

func (b *Task) TableName() string {
	return "ToDo"
}