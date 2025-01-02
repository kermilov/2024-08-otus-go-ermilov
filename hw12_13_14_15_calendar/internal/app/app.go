package app

import (
	"context"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
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

func New(logger Logger, storage Storage) *App { //nolint:revive
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error { //nolint:revive
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
