package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/util"
)

// Общий интерфейс логгера на разные реализации.
type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type ListenerFunc func(context.Context, Application, []byte) error

type Consumer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Application interface {
	SaveNotification(ctx context.Context, id, title string, datetime time.Time, userid int64) error
}

func SaveNotification(ctx context.Context, app Application, bytes []byte) error {
	notification := util.Notification{}
	err := json.Unmarshal(bytes, &notification)
	if err != nil {
		return err
	}
	err = app.SaveNotification(ctx,
		notification.ID,
		notification.Title,
		notification.DateTime,
		notification.UserID,
	)
	if err != nil {
		return err
	}
	return nil
}
