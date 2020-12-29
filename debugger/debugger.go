package debugger

import (
	"encoding/json"
	"fmt"
	"github.com/go-delve/delve/service"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/go-delve/delve/service/rpccommon"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

// /go-delve/delve@v1.5.1/pkg/gobuild/gobuild.go
func GoBuild(debugName string) error {
	args := []string{"-o", debugName}
	return command("build", args...)
}

func command(command string, args ...string) error {
	allArgs := []string{command}
	allArgs = append(allArgs, args...)
	goBuild := exec.Command("go", allArgs...)
	goBuild.Stderr = os.Stderr
	goBuild.Stdout = os.Stdout
	return goBuild.Run()
}

func Start() error {
	addr := "127.0.0.1:9229"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	debugName := "sample"
	processArgs := []string{debugName}

	server := rpccommon.NewServer(&service.Config{
		Listener:           listener,
		ProcessArgs:        processArgs,
		AcceptMulti:        false,
		APIVersion:         2,
		CheckLocalConnUser: false,
		Debugger: debugger.Config{
			WorkingDir:     "",
			Backend:        "default",
			CheckGoVersion: false,
		},
	})

	if err := server.Run(); err != nil {
		return err
	}

	path, err := filepath.Abs("sample.go")
	if err != nil {
		return err
	}

	fmt.Println(path)

	client := rpc2.NewClient(addr)
	_, err = client.CreateBreakpoint(&api.Breakpoint{
		Line: 6,
		File: path,
	})
	if err != nil {
		return err
	}

	state, err := client.GetState()
	if err != nil {
		return err
	}

	variables, err := client.ListLocalVariables(
		api.EvalScope{
			GoroutineID:  state.CurrentThread.GoroutineID,
			Frame:        0,
			DeferredCall: 0,
		},
		api.LoadConfig{},
	)
	if err != nil {
		return err
	}
	bytes, _ := json.Marshal(variables)
	fmt.Printf("Thread: %+v\n", string(bytes))

	//stackTraceOpts := api.StacktraceOptions(0)
	//loadConfig := &api.LoadConfig{}
	//stacktrace, err := client.Stacktrace(state.CurrentThread.GoroutineID, 10, stackTraceOpts, loadConfig)
	//if err != nil {
	//	return err
	//}
	//
	//marshal, err := json.Marshal(stacktrace)
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Printf("Stack trance: %+v\n", string(marshal))

	return nil
}
