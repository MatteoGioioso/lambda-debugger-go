package main

import (
	"fmt"
	"lambda-debugger-go/collector"
	"lambda-debugger-go/debugger"
	"lambda-debugger-go/ipc"
	"lambda-debugger-go/logger"
	"lambda-debugger-go/utils"
	"os"
	"strconv"
	"strings"
)

var (
	rawPid        = os.Getenv("DEBUG_TARGET_PID")
	ipcClientName = os.Getenv("DEBUG_NAMED_PIPE")
	deb           = debugger.New("handler.go", "handler")
	log           = logger.New(true, true)
	ipcClient     = ipc.New(ipcClientName)
	steps         = debugger.NewSteps()
	files         = debugger.NewFiles()
	col           = collector.New()
)

func recordExecution(deb debugger.Port, goroutineID int) {
	currentGoroutineID := goroutineID

	for {
		variables, err := deb.GetLocalVariables(currentGoroutineID)
		if err != nil {
			fmt.Println(err)
			return
		}
		stackTrace, err := deb.GetStackTrace(currentGoroutineID)
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.Contains(stackTrace[0].File, "lambda-debugger-go") {
			steps.AddStep(variables, stackTrace)
			if err := files.Add(stackTrace[0].File); err != nil {
				fmt.Println(err)
				return
			}

			state, err := deb.StepIn()
			if err != nil {
				fmt.Println(err)
				return
			}

			goroutineID = state.CurrentThread.GoroutineID
			continue
		} else {
			state, err := deb.StepOut()
			if err != nil {
				fmt.Println(err)
				return
			}

			// If no thread is present means that the execution has been finished
			if state.CurrentThread != nil {
				goroutineID = state.CurrentThread.GoroutineID
			}
			continue
		}
	}
}

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

	if err := deb.AddBreakpoint(26); err != nil {
		fmt.Println(err)
		return
	}
	log.Info("Breakpoint created")

	state := deb.Continue()
	recordExecution(deb, state.CurrentThread.GoroutineID)

	stepsDTO := debugger.ToStepsDTO(*steps)
	filesDTO := debugger.ToFilesDTO(*files)
	if err := col.InjectDebuggerOutputIntoHtml(stepsDTO, filesDTO); err != nil {
		fmt.Println(err)
		return
	}
}
