package main

import (
	"todoapp/tasks"
	"todoapp/users"
	"todoapp/auth"
)


func main() {

	tasks.SetupDB()
	users.SetupUDB()

	r := auth.SetupRouter()
	users.InitUsers(r)
	tasks.InitTasks(r)
	
	r.Run(":8080") // Main File 
}
