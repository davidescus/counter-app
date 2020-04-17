package main

import (
	"counter-app/pkg/app"
	"fmt"
)

// TODO add WebServer
// TODO add logging, access log, errors, warnings, etc
// TODO graceful shutdown
// TODO create swagger, add endpoint to WebServer
// TODO make it persistent
// TODO scale it with multiple instances

func main() {
	fmt.Println(" --- start --- ")

	a := app.New()

	// TODO maybe trigger this when create app
	a.Start()

	err := a.Store("one tho one three")
	// TODO maybe not return error, writing to memory
	_ = err



	fmt.Println("One: ", a.Get([]string{"one"}))
	fmt.Println("Three: ", a.Get([]string{"three", "one"}))
	fmt.Println("Four: ", a.Get([]string{"four"}))
}
