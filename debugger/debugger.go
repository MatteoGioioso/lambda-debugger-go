package debugger

import (
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

type Variable struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type deb struct {
	client         service.Client
	server         service.Server
	address        string
	binaryFileName string
	sourceFileName string
}

func New(source, binary string) *deb {
	return &deb{
		binaryFileName: binary,
		sourceFileName: source,
		address:        "127.0.0.1:9229",
	}
}

func (d *deb) InitServer() error {
	listener, err := net.Listen("tcp", d.address)
	if err != nil {
		return err
	}

	processArgs := []string{d.binaryFileName}

	server := rpccommon.NewServer(&service.Config{
		Listener:           listener,
		ProcessArgs:        processArgs,
		AcceptMulti:        false,
		APIVersion:         2,
		CheckLocalConnUser: false,
		Debugger: debugger.Config{
			Backend:        "default",
			CheckGoVersion: false,
		},
	})

	if err := server.Run(); err != nil {
		return err
	}

	return nil
}

func (d *deb) InitClient() error {
	client := rpc2.NewClient(d.address)
	d.client = client
	return nil
}

func (d deb) AddBreakpoint(line int) error {
	path, err := filepath.Abs(d.sourceFileName)
	if err != nil {
		return err
	}
	_, err = d.client.CreateBreakpoint(&api.Breakpoint{
		Line: line,
		File: path,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d deb) GetClient() service.Client {
	return d.client
}

func (d deb) GetLocalVariables(goRoutineID int) ([]Variable, error) {
	variables, err := d.client.ListLocalVariables(
		api.EvalScope{
			GoroutineID: goRoutineID,
			Frame:       0,
		},
		api.LoadConfig{
			FollowPointers:     true,
			MaxStructFields:    -1,
			MaxVariableRecurse: 1,
			MaxStringLen:       100,
			MaxArrayValues:     100,
		},
	)

	vars := make([]Variable, 0)
	for _, variable := range variables {
		vars = append(vars, Variable{
			Name:  variable.Name,
			Type:  variable.RealType,
			Value: variable.Value,
		})
	}

	return vars, err
}

func (d deb) GetStackTrace(goRoutineID int) error {
	stackTraceOpts := api.StacktraceOptions(0)
	loadConfig := &api.LoadConfig{}
	stacktrace, err := d.client.Stacktrace(goRoutineID, 50, stackTraceOpts, loadConfig)
	if err != nil {
		return err
	}

	fmt.Println(stacktrace)

	return nil
}

func (d deb) Continue() *api.DebuggerState {
	state := <-d.client.Continue()
	return state
}

func (d deb) Step() (*api.DebuggerState, error) {
	step, err := d.client.Step()
	if err != nil {
		return &api.DebuggerState{}, err
	}

	return step, err
}
