package logger

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	errorLvl  = "ERROR"
	normalLvl = "NORMAL"
)

type logger struct {
	level      string
	pid        int
	isDebugger bool
	debug      bool
	message    interface{}
	enc        *json.Encoder
}

func New(debug bool, isDebugger bool) *logger {
	enc := json.NewEncoder(os.Stdout)
	return &logger{
		enc:        enc,
		debug:      debug,
		isDebugger: isDebugger,
	}
}

func (l logger) setProcessType(fullMsg string) string {
	if l.isDebugger {
		fullMsg = fullMsg + "(Debugger): "
	} else {
		fullMsg = fullMsg + "(Handler): "
	}

	return fullMsg
}

func (l logger) Info(message string) {
	fullMsg := fmt.Sprintf("[Lambda-debugger - %v]: ", normalLvl)
	fullMsg = l.setProcessType(fullMsg)
	fullMsg = fullMsg + message

	if l.debug {
		l.enc.Encode(fullMsg)
	}
}

func (l logger) Failure(err error) {
	fullMsg := fmt.Sprintf("[Lambda-debugger - %v]: ", errorLvl)
	fullMsg = l.setProcessType(fullMsg)
	fullMsg = fullMsg + err.Error()

	if l.debug {
		l.enc.Encode(fullMsg)
	}
}
