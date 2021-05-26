package main

import (
	"todoapp/router"
	"todoapp/config"
)

var err error

func main() {

	config.SetupDB()

	r := router.SetupRouter()
	
	r.Run()

}
