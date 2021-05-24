package main

import (
	"todoapp/router"
)

var err error

func main() {
	
	r := router.SetupRouter()
	r.Run()
	
}
