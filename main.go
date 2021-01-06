package main

import (
	"fmt"
	"lambda-debugger-go/debugger"
)

var deb = debugger.New("sample.go", "sample")

func main() {
	if err := deb.InitServer(); err != nil {
		fmt.Println(err)
	}

	if err := deb.InitClient(); err != nil {
		fmt.Println(err)
	}

	if err := deb.AddBreakpoint(6); err != nil {
		fmt.Println(err)
	}

	state := deb.Continue()
	variables, err := deb.GetLocalVariables(state.CurrentThread.GoroutineID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(variables)
}
