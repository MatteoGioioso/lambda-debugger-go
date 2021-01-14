package utils

import (
	"os"
	"strconv"
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
