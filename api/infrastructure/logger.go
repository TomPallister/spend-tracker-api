package infrastructure

import "fmt"

// Logger ...
type Logger interface {
	Info(message string)
	Error(message string, err error)
}

// ConsoleLogger ...
type ConsoleLogger struct {
}

// Info ...
func (logger ConsoleLogger) Info(message string) {
	fmt.Println(message)
}

func (logger ConsoleLogger) Error(message string, err error) {
	fmt.Println(message, err)
}

// NilLogger ...
type NilLogger struct {
}

// Info ...
func (logger NilLogger) Info(message string) {
}

func (logger NilLogger) Error(message string, err error) {
}
