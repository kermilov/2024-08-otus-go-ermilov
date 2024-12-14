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
	Create(event storage.Event) (storage.Event, error)
	// Обновить (ID события, событие);
	Update(id string, event storage.Event) error
	// Удалить (ID события);
	Delete(id string) error
	// СписокСобытийНаДень (дата);
	FindByDay(date time.Time) ([]storage.Event, error)
	// СписокСобытийНаНеделю (дата начала недели);
	FindByWeek(date time.Time) ([]storage.Event, error)
	// СписокСобытийНaМесяц (дата начала месяца).
	FindByMonth(date time.Time) ([]storage.Event, error)
	// пр. на усмотрение разработчика.
	FindByID(id string) (storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
