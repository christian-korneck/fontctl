package main

import (
	"fmt"
	"os"
)

type debugLogger interface {
	Error(message string)
	Warn(message string)
	Info(message string)
}

type CliDebugLogger struct{}

func (CliDebugLogger) Error(message string) {
	fmt.Fprintf(os.Stderr, "[DEBUG][ERR] - %s\n", message)
}

func (CliDebugLogger) Warn(message string) {
	fmt.Fprintf(os.Stderr, "[DEBUG][WRN] - %s\n", message)
}

func (CliDebugLogger) Info(message string) {
	fmt.Fprintf(os.Stdout, "[DEBUG][INF] - %s\n", message)
}
