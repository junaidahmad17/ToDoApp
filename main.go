package main

import (
	"todoapp/auth"
	"todoapp/tasks"
	"todoapp/users"
)

func main() {

	tasks.SetupDB()
	users.SetupUDB()

	r := auth.SetupRouter()
	users.InitUsers(r)
	tasks.InitTasks(r)

	r.Run(":8080")
}
