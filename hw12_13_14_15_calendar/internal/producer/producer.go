package producer

import (
	"context"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
)

// Общий интерфейс логгера на разные реализации планировщика.
type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

// Общий интерфейс приложения на разные реализации планировщика.
type Producer interface {
	Start(ctx context.Context) error
	ScheduledProcess(ctx context.Context)
	Stop(ctx context.Context) error
}

type Application interface {
	FindForSendNotification(ctx context.Context, date time.Time) ([]storage.Event, error)
	SetIsSendNotification(ctx context.Context, ids []string) error
	DeleteOldEvents(ctx context.Context, date time.Time) error
}
