package app

import (
	"context"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	// Создать (событие);
	Create(ctx context.Context, event storage.Event) (storage.Event, error)
	// Обновить (ID события, событие);
	Update(ctx context.Context, id string, event storage.Event) error
	// Удалить (ID события);
	Delete(ctx context.Context, id string) error
	// СписокСобытийНаДень (дата);
	FindByDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	// СписокСобытийНаНеделю (дата начала недели);
	FindByWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	// СписокСобытийНaМесяц (дата начала месяца).
	FindByMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
	// пр. на усмотрение разработчика.
	FindByID(ctx context.Context, id string) (storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		storage: storage,
		logger:  logger,
	}
}

func (a *App) CreateEvent(ctx context.Context,
	id string,
	title string,
	datetime time.Time,
	duration *time.Duration,
	userid int64,
) (
	*storage.Event, error,
) {
	a.logger.Info("create event")
	eventDuration := time.Hour
	if duration != nil {
		eventDuration = *duration
	}
	event, err := a.storage.Create(ctx,
		storage.Event{
			ID:       id,
			Title:    title,
			DateTime: datetime,
			Duration: eventDuration,
			UserID:   userid,
		},
	)
	if err != nil {
		a.logger.Error(err.Error())
	}
	return &event, err
}

func (a *App) UpdateEvent(ctx context.Context,
	id string,
	title string,
	datetime time.Time,
	duration *time.Duration,
	userid int64,
) error {
	a.logger.Info("update event")
	eventDuration := time.Hour
	if duration != nil {
		eventDuration = *duration
	}
	err := a.storage.Update(ctx,
		id,
		storage.Event{
			Title:    title,
			DateTime: datetime,
			Duration: eventDuration,
			UserID:   userid,
		})
	if err != nil {
		a.logger.Error(err.Error())
	}
	return err
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	a.logger.Info("delete event")
	err := a.storage.Delete(ctx, id)
	if err != nil {
		a.logger.Error(err.Error())
	}
	return err
}

func (a *App) FindEventByID(ctx context.Context, id string) (storage.Event, error) {
	a.logger.Info("find event by id")
	event, err := a.storage.FindByID(ctx, id)
	if err != nil {
		a.logger.Error(err.Error())
	}
	return event, err
}

func (a *App) FindEventByDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Info("find event by day")
	events, err := a.storage.FindByDay(ctx, date)
	if err != nil {
		a.logger.Error(err.Error())
	}
	return events, err
}

func (a *App) FindEventByWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Info("find event by week")
	events, err := a.storage.FindByWeek(ctx, date)
	if err != nil {
		a.logger.Error(err.Error())
	}
	return events, err
}

func (a *App) FindEventByMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Info("find event by month")
	events, err := a.storage.FindByMonth(ctx, date)
	if err != nil {
		a.logger.Error(err.Error())
	}
	return events, err
}
