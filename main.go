package main

import (
	"log"
	"os"

	"todoapp/auth"
	"todoapp/tasks"
	"todoapp/users"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	tasks.SetupDB()
	users.SetupUDB()
}
// redis 
func main() {
	tasks.Remind()

	r := auth.SetupRouter()
	users.InitUsers(r)
	tasks.InitTasks(r)
	r.Run(os.Getenv("PORT"))
	
}

