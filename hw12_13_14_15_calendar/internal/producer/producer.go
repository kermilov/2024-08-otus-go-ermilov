package producer

// Общий интерфейс логгера на разные реализации планировщика.
type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}
