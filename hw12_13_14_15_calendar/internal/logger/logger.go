package logger

import "fmt"

type LogLevel int

const (
	Error LogLevel = iota
	Warning
	Info
	Debug
)

var logLevels = map[string]LogLevel{
	"ERROR":   Error,
	"WARNING": Warning,
	"INFO":    Info,
	"DEBUG":   Debug,
}

type Logger struct {
	logLevel LogLevel
}

func New(level string) *Logger {
	logLevel, isOk := logLevels[level]
	if !isOk {
		panic("неизвестный уровень логирования")
	}
	return &Logger{logLevel}
}

func (l Logger) Error(msg string) {
	if l.logLevel >= Error {
		fmt.Println(msg)
	}
}

func (l Logger) Warning(msg string) {
	if l.logLevel >= Warning {
		fmt.Println(msg)
	}
}

func (l Logger) Info(msg string) {
	if l.logLevel >= Info {
		fmt.Println(msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.logLevel >= Debug {
		fmt.Println(msg)
	}
}
