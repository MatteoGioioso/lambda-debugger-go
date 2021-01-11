package main

import (
	"fmt"
	"lambda-debugger-go/debugger"
	"lambda-debugger-go/ipc"
	"lambda-debugger-go/logger"
	"lambda-debugger-go/utils"
	"os"
	"strconv"
)

var (
	rawPid        = os.Getenv("DEBUG_TARGET_PID")
	ipcClientName = os.Getenv("DEBUG_NAMED_PIPE")
	deb           = debugger.New("handler.go", "handler")
	log           = logger.New(true, true)
	ipcClient     = ipc.New(ipcClientName)
)

func main() {
	pid, err := strconv.Atoi(rawPid)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := deb.InitServer(pid); err != nil {
		fmt.Println(err)
		return
	}
	log.Info("Server created")

	if err := deb.InitClient(); err != nil {
		fmt.Println(err)
		return
	}
	log.Info("Client created")

	callback := func() {
		if err := deb.Clean(); err != nil {
			fmt.Println(err)
		}
		if err := ipcClient.Close(); err != nil {
			fmt.Println(err)
		}
	}

	defer utils.OnPanicOrExit(callback)
	utils.OnSignTerm(callback)

	if err := deb.AddBreakpoint(10); err != nil {
		fmt.Println(err)
		return
	}
	log.Info("Breakpoint created")

	state := deb.Continue()

	variables, err := deb.GetLocalVariables(state.CurrentThread.GoroutineID)
	if err != nil {
		fmt.Println(err)
		return
	}
	stackTrace, err := deb.GetStackTrace(state.CurrentThread.GoroutineID)
	if err != nil {
		fmt.Println(err)
		return
	}

	steps := debugger.NewSteps()
	steps.AddStep(variables, stackTrace)
	fmt.Printf("%+v\n", steps)
}
