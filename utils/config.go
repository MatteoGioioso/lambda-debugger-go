package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

const (
	lambdaDebuggerNameSpace = "lambda-debugger"
)

func GetMaxArrayValues() (int, error) {
	env := os.Getenv("LAMBDA_DEBUGGER_MAX_ARRAY_VALUES")
	length, err := strconv.Atoi(env)
	if err != nil {
		return 0, err
	}

	return length, err
}

func SetMaxArrayValues() (int, error) {
	env := os.Getenv("LAMBDA_DEBUGGER_MAX_ARRAY_VALUES")
	if env == "" {
		return 100, nil
	}

	length, err := strconv.Atoi(env)
	if err != nil {
		return 0, err
	}

	return length, err
}

func SetOutputPath() string {
	env := os.Getenv("LAMBDA_DEBUGGER_OUTPUT_PATH")
	if env == "" {
		return filepath.Join("/tmp", lambdaDebuggerNameSpace)
	}

	return env
}

func SetFilePath() string {
	env := os.Getenv("LAMBDA_DEBUGGER_FILE_PATH")
	if env == "" {
		return os.Getenv("LAMBDA_TASK_ROOT")
	}

	return env
}
