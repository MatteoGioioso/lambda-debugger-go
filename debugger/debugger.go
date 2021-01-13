package debugger

import (
	"github.com/go-delve/delve/service"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/debugger"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/go-delve/delve/service/rpccommon"
	"net"
	"path/filepath"
)

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

func (d *deb) InitServer(pid int) error {
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
			AttachPid: pid,
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

func (d deb) GetLocalVariables(goRoutineID int) ([]api.Variable, error) {
	variables, err := d.client.ListLocalVariables(
		api.EvalScope{
			GoroutineID: goRoutineID,
			Frame:       0,
		},
		api.LoadConfig{
			FollowPointers:     true,
			MaxStructFields:    -1,
			MaxVariableRecurse: 5,
			MaxStringLen:       100,
			MaxArrayValues:     100,
		},
	)

	return variables, err
}

func (d deb) GetStackTrace(goRoutineID int) ([]api.Stackframe, error) {
	stackTraceOpts := api.StacktraceOptions(0)
	loadConfig := &api.LoadConfig{}
	stacktrace, err := d.client.Stacktrace(goRoutineID, 50, stackTraceOpts, loadConfig)
	if err != nil {
		return nil, err
	}

	return stacktrace, err
}

func (d deb) Continue() *api.DebuggerState {
	state := <-d.client.Continue()
	return state
}

func (d deb) StepIn() (*api.DebuggerState, error) {
	step, err := d.client.Step()
	if err != nil {
		return &api.DebuggerState{}, err
	}

	return step, err
}

func (d deb) StepOut() (*api.DebuggerState, error) {
	return d.client.StepOut()
}

func (d deb) Clean() error {
	if err := d.client.Detach(true); err != nil {
		return err
	}

	if d.server != nil {
		if err := d.server.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (d deb) GetState() (*api.DebuggerState, error) {
	state, err := d.client.GetState()
	if err != nil {
		return nil, err
	}

	return state, nil
}

type Port interface {
	GetLocalVariables(goRoutineID int) ([]api.Variable, error)
	GetStackTrace(goRoutineID int) ([]api.Stackframe, error)
	StepIn() (*api.DebuggerState, error)
	StepOut() (*api.DebuggerState, error)
}
