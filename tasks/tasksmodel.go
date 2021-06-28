package tasks

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

type Task struct {
	ID          uint      `json:"idt"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Create_DT   time.Time `json:"create_DT"`
	Due_DT      time.Time `json:"due_DT"`
	Com_status  bool      `json:"com_status"`
	Com_DT      time.Time `json:"com_DT"`
	Attachment  string
	Uid         int
}

type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Com_status  bool   `json:"com_status"`
}

func SetupDB() {
	flag := false

	if DB == nil {

		if flag {
			e := os.Remove("ToDo.db")
			if e != nil {
				log.Fatalln(e)
			}
		}
		DB, err = gorm.Open(sqlite.Open(os.Getenv("DBADD")+"ToDo.db"), &gorm.Config{})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// form data
