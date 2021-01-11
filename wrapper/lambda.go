package wrapper

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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
	return func() {
		if err := ipcClient.Create(); err != nil {
			fmt.Println("ipcClient creation failed: ", err)
		}
		pid := os.Getpid()
		if err := utils.GoBuild("main"); err != nil {
			fmt.Println("Go build failed: ", err)
		}

		cmd := exec.Command("./main")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("DEBUG_TARGET_PID=%v", pid),
			fmt.Sprintf("DEBUG_NAMED_PIPE=%v", ipcClient.GetName()),
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

		f, ok := handler.(func(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error))
		if !ok {
			fmt.Println("NOT a function")
		}
		f(events.APIGatewayProxyRequest{})
	}
}
