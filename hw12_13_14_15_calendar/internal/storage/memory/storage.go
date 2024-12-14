package memorystorage

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	storage map[string]storage.Event
	mu      sync.RWMutex
}

func New() *Storage {
	return &Storage{storage: make(map[string]storage.Event)}
}

// Создать (событие).
func (s *Storage) Create(event storage.Event) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	s.storage[event.ID] = event
	return event, nil
}

// Обновить (ID события, событие).
func (s *Storage) Update(id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[id] = event
	return nil
}

// Удалить (ID события).
func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.storage, id)
	return nil
}

// СписокСобытийНаДень (дата).
func (s *Storage) FindByDay(date time.Time) ([]storage.Event, error) {
	startOfDay := storage.GetStartOfDay(date)
	endOfDay := storage.GetEndOfDay(date)
	return s.findByDateTimeBetween(startOfDay, endOfDay)
}

// СписокСобытийНаНеделю (дата начала недели).
func (s *Storage) FindByWeek(date time.Time) ([]storage.Event, error) {
	startOfWeek := storage.GetStartOfWeek(date)
	endOfWeek := storage.GetEndOfWeek(date)
	return s.findByDateTimeBetween(startOfWeek, endOfWeek)
}

// СписокСобытийНaМесяц (дата начала месяца).
func (s *Storage) FindByMonth(date time.Time) ([]storage.Event, error) {
	startOfMonth := storage.GetStartOfMonth(date)
	endOfMonth := storage.GetEndOfMonth(date)
	return s.findByDateTimeBetween(startOfMonth, endOfMonth)
}

// пр. на усмотрение разработчика.
func (s *Storage) FindByID(id string) (storage.Event, error) {
	result, isOk := s.storage[id]
	if !isOk {
		return storage.Event{}, storage.ErrEventNotFound
	}
	return result, nil
}

func (s *Storage) findByDateTimeBetween(startDate time.Time, endDate time.Time) ([]storage.Event, error) {
	result := make([]storage.Event, 0)
	for _, v := range s.storage {
		if (v.DateTime.Compare(startDate) >= 0) && (v.DateTime.Compare(endDate) <= 0) {
			result = append(result, v)
		}
	}
	return result, nil
}
