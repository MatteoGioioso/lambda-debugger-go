package main

import (
	"fmt"
	"lambda-debugger-go/ipc"
	"lambda-debugger-go/utils"
	"os"
)

//var deb = debugger.New("handler.go", "handler")

func main() {
	fmt.Println("Started...")
	//rawPid := os.Getenv("DEBUG_TARGET_PID")
	ipcClientName := os.Getenv("DEBUG_NAMED_PIPE")

	ipcClient := ipc.New(ipcClientName)

	//pid, err := strconv.Atoi(rawPid)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//if err := deb.InitServer(pid); err != nil {
	//	fmt.Println(err)
	//}
	//
	//if err := deb.InitClient(); err != nil {
	//	fmt.Println(err)
	//}

	if err := ipcClient.Send("Done"); err != nil {
		fmt.Println(err)
	}

	callback := func() {
		//if err := deb.Clean(); err != nil {
		//	fmt.Println(err)
		//}
		if err := ipcClient.Close(); err != nil {
			fmt.Println(err)
		}
	}

	utils.OnPanicOrExit(callback)
	utils.OnSignTerm(callback)

	//if err := deb.AddBreakpoint(19); err != nil {
	//	fmt.Println(err)
	//}

	//state := deb.Continue()
	//variables, err := deb.GetLocalVariables(state.CurrentThread.GoroutineID)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(variables)
}
