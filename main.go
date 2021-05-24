package main

import (
	"todoapp/router"
)

var err error

func main() {
	// handler.SetupDB()
	r := router.SetupRouter()
	r.Run()
	// Comment 
}
