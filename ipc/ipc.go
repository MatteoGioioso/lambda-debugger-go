package ipc

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"os"
)

// ipc inter-process communication
// this package is responsible for communicating between the main process and the debugger
// it will mainly signal when the debugger is ready and the process can continue
type ipc struct {
	name string
}

func New(name string) *ipc {
	return &ipc{name: name}
}

func (i ipc) GetName() string {
	return i.name
}

func (i ipc) Create() error  {
	if err := unix.Mkfifo(i.name, 0666); err != nil {
		return err
	}

	return nil
}

func (i ipc) WaitForMessage() (string, error) {
	file, err := os.OpenFile(i.name, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	fmt.Println("Waiting for someone to write something")
	if _, err := io.Copy(&buff, file); err != nil {
		return "", err
	}
	defer file.Close()

	return buff.String(), err
}

func (i ipc) Send(msg string) error {
	stdout, err := os.OpenFile(i.name, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	if _, err := stdout.Write([]byte(msg)); err != nil {
		return err
	}
	if err := stdout.Close(); err != nil {
		return err
	}

	return nil
}

func (i ipc) Close() error {
	if err := os.Remove(i.name); err != nil {
		return err
	}

	return nil
}
