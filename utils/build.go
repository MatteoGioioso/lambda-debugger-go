package utils

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func GoBuild(forkProcBinName string) error {
	args := []string{"-o", forkProcBinName, "../main.go"}
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

func OnPanicOrExit(callback func()) {
	defer func() {
		r := recover()

		callback()

		if r != nil {
			panic(r)
		}
	}()
}

func OnSignTerm(callback func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		sig := <-sigs

		// Log the received signal
		fmt.Printf("LOG: Caught sig ")
		fmt.Println(sig)

		callback()

		// Gracefully exit.
		// (Use runtime.GoExit() if you need to call defers)
		os.Exit(0)
	}()
}
