package wrapper

import (
	"fmt"
	"lambda-debugger-go/ipc"
	"lambda-debugger-go/logger"
	"lambda-debugger-go/utils"
	"os"
	"os/exec"
)

var (
	log = logger.New(true, false)
	ipcClient = ipc.New("debugger-pipe")
)

func LambdaWrapper(handler interface{}) interface{} {
	return func() () {
		pid := os.Getpid()
		cmd := exec.Command("./main")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("DEBUG_TARGET_PID=%v", pid),
			fmt.Sprintf("DEBUG_NAMED_PIPE=%v", ipcClient.GetName()),
			fmt.Sprintf("LAMBDA_DEBUGGER_MAX_ARRAY_VALUES=100"),
			fmt.Sprintf("LAMBDA_DEBUGGER_OUTPUT_PATH=%v", utils.SetOutputPath()),
			fmt.Sprintf("LAMBDA_DEBUGGER_FILE_PATH=%v", utils.SetFilePath()),
		)
		if err := cmd.Start(); err != nil {
			fmt.Println("Command start failed: ", err)
		}

		log.Info("Waiting for debugger")
		message, err := ipcClient.WaitForMessage()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(message)

		f := handler.(func())
		f()
		//f, ok := handler.(func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error))
		//if !ok {
		//	fmt.Println("NOT a function")
		//}
		//resp, err := f(events.APIGatewayProxyRequest{})
		//return resp, err
	}
}
