package main

import (
	"counter-app/pkg/app"
	"fmt"
)

func main() {
	fmt.Println(" --- start --- ")

	a := app.New()

	// TODO maybe trigger this when create app
	a.Start()

	err := a.Store("one tho one three")
	// TODO maybe not return error, writing to memory
	_ = err



	fmt.Println("One: ", a.Get("one"))
	fmt.Println("Three: ", a.Get("three"))
	fmt.Println("Four: ", a.Get("four"))
}
