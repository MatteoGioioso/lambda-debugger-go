package main

import (
	"fmt"
	"lambda-debugger-go/debugger"
)

func main() {
	if err := debugger.Start(); err != nil {
		fmt.Println("Error: ", err)
	}
}
