package server

import (
	"context"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
)

// Указание формата времени.
const Layout = time.RFC3339

// Общий интерфейс логгера на разные реализации сервера.
type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

// Общий интерфейс приложения на разные реализации сервера.
type Application interface {
	CreateEvent(ctx context.Context,
		id string, title string, datetime time.Time, duration *time.Duration, userid int64,
		notificationDuration *time.Duration,
	) (
		*storage.Event, error,
	)
	UpdateEvent(ctx context.Context,
		id string, title string, datetime time.Time, duration *time.Duration, userid int64,
		notificationDuration *time.Duration) error
	DeleteEvent(ctx context.Context, id string) error
	FindEventByDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventByWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventByMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
	FindEventByID(ctx context.Context, id string) (storage.Event, error)
}
