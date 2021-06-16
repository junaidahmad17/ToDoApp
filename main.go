package main

import (
	"os"
	
	"todoapp/auth"
	"todoapp/tasks"
	"todoapp/users"

	"github.com/joho/godotenv"
)
func init() {
	_ = godotenv.Load()
	tasks.SetupDB()
	users.SetupUDB()
}
func main() {
	
	r := auth.SetupRouter()
	users.InitUsers(r)
	tasks.InitTasks(r)
	r.Run(os.Getenv("PORT"))
}
