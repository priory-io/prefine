package logging

import (
	"fmt"
	"log"
	"os"
)

type ConsoleLogger struct {
	verbose bool
	logger  *log.Logger
}

func NewConsoleLogger(verbose bool) *ConsoleLogger {
	return &ConsoleLogger{
		verbose: verbose,
		logger:  log.New(os.Stdout, "", 0),
	}
}

func (l *ConsoleLogger) Info(msg string, args ...any) {
	l.logger.Printf(msg, args...)
}

func (l *ConsoleLogger) Error(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+msg+"\n", args...)
}

func (l *ConsoleLogger) Debug(msg string, args ...any) {
	if l.verbose {
		l.logger.Printf("Debug: "+msg, args...)
	}
}
