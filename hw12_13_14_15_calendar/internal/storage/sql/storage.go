package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/migrations"
	goose "github.com/pressly/goose/v3"
)

type Storage struct {
	db   *sql.DB
	conn *sql.Conn
}

func New(dsn string) *Storage {
	db, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		panic(fmt.Errorf("не удалось установить первоначальное соединение с базой данных: %w", err))
	}

	goose.SetBaseFS(migrations.EmbedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
	return &Storage{db: db}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("не удалось установить соединение с базой данных: %w", err)
	}
	s.conn = conn
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	err := s.conn.Close()
	if err != nil {
		return fmt.Errorf("не удалось закрыть соединение с базой данных: %w", err)
	}
	return nil
}

const (
	checkDateBusyQuery = `select 1 from calendar.event where $1 between datetime and datetime + duration 
	or ($1::timestamp + $2::interval) between datetime and datetime + duration`
	insertQuery = `insert into event (id, title, datetime, duration, userid) 
	               values ($1, $2, $3, $4, $5) returning id, title, datetime, duration, userid`
	updateQuery                = `update event set title = $1, datetime = $2, duration = $3, userid = $4 where id = $5`
	deleteQuery                = `delete from event where id = $1`
	findByDateTimeBetweenQuery = `select id, title, datetime, duration, userid from event where datetime between $1 and $2`
	findByIDQuery              = `select id, title, datetime, duration, userid from event where id = $1`
)

// Создать (событие).
func (s *Storage) Create(ctx context.Context, event storage.Event) (storage.Event, error) {
	err := s.Connect(ctx)
	if err != nil {
		return storage.Event{}, err
	}
	defer s.Close(ctx)
	err = s.checkDateBusy(ctx, event)
	if err != nil {
		return storage.Event{}, err
	}
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	_, err = s.conn.ExecContext(
		ctx, insertQuery, event.ID, event.Title, event.DateTime, fmt.Sprintf("%v", event.Duration), event.UserID,
	)
	if err != nil {
		return storage.Event{}, fmt.Errorf("не удалось создать событие: %w", err)
	}
	return storage.Event{}, nil
}

// Обновить (ID события, событие).
func (s *Storage) Update(ctx context.Context, id string, event storage.Event) error {
	err := s.Connect(ctx)
	if err != nil {
		return err
	}
	defer s.Close(ctx)
	err = s.checkDateBusy(ctx, event)
	if err != nil {
		return err
	}
	_, err = s.conn.ExecContext(
		ctx, updateQuery, event.Title, event.DateTime, fmt.Sprintf("%v", event.Duration), event.UserID, id,
	)
	if err != nil {
		return fmt.Errorf("не удалось обновить событие: %w", err)
	}
	return nil
}

// Удалить (ID события).
func (s *Storage) Delete(ctx context.Context, id string) error {
	err := s.Connect(ctx)
	if err != nil {
		return err
	}
	defer s.Close(ctx)
	_, err = s.conn.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить событие: %w", err)
	}
	return nil
}

// СписокСобытийНаДень (дата).
func (s *Storage) FindByDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	err := s.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer s.Close(ctx)
	startOfDay := storage.GetStartOfDay(date)
	endOfDay := storage.GetEndOfDay(date)
	return s.findByDateTimeBetween(ctx, startOfDay, endOfDay)
}

// СписокСобытийНаНеделю (дата начала недели).
func (s *Storage) FindByWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	err := s.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer s.Close(ctx)
	startOfWeek := storage.GetStartOfWeek(date)
	endOfWeek := storage.GetEndOfWeek(date)
	return s.findByDateTimeBetween(ctx, startOfWeek, endOfWeek)
}

// СписокСобытийНaМесяц (дата начала месяца).
func (s *Storage) FindByMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	err := s.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer s.Close(ctx)
	startOfMonth := storage.GetStartOfMonth(date)
	endOfMonth := storage.GetEndOfMonth(date)
	return s.findByDateTimeBetween(ctx, startOfMonth, endOfMonth)
}

// пр. на усмотрение разработчика.
func (s *Storage) FindByID(ctx context.Context, id string) (storage.Event, error) {
	err := s.Connect(ctx)
	if err != nil {
		return storage.Event{}, err
	}
	defer s.Close(ctx)
	row := s.conn.QueryRowContext(ctx, findByIDQuery, id)
	return s.rowToEvent(row)
}

func (s *Storage) checkDateBusy(ctx context.Context, event storage.Event) error {
	result, err := s.conn.QueryContext(
		ctx, checkDateBusyQuery, event.DateTime, fmt.Sprintf("%v", event.Duration),
	)
	if err != nil {
		return fmt.Errorf("не удалось проверить свободность даты: %w", err)
	}
	defer result.Close()
	if result.Next() {
		return storage.ErrDateBusy
	}
	return nil
}

func (s *Storage) findByDateTimeBetween(ctx context.Context, startDate time.Time, endDate time.Time) (
	[]storage.Event, error,
) {
	rows, err := s.conn.QueryContext(ctx, findByDateTimeBetweenQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить события: %w", err)
	}
	defer rows.Close()
	return s.rowsToEvents(rows)
}

func (s *Storage) rowsToEvents(rows *sql.Rows) ([]storage.Event, error) {
	result := make([]storage.Event, 0)
	for rows.Next() {
		var event storage.Event
		var duration string
		err := rows.Scan(&event.ID, &event.Title, &event.DateTime, &duration, &event.UserID)
		if err != nil {
			return nil, fmt.Errorf("не удалось сканировать событие: %w", err)
		}
		event.Duration, err = parsePostgresInterval(duration)
		if err != nil {
			return nil, fmt.Errorf("не удалось преобразовать duration: %w", err)
		}
		result = append(result, event)
	}
	return result, nil
}

func (s *Storage) rowToEvent(row *sql.Row) (storage.Event, error) {
	var event storage.Event
	var duration string
	err := row.Scan(&event.ID, &event.Title, &event.DateTime, &duration, &event.UserID)
	if err != nil {
		return storage.Event{}, fmt.Errorf("не удалось сканировать событие: %w", err)
	}
	event.Duration, err = parsePostgresInterval(duration)
	if err != nil {
		return storage.Event{}, fmt.Errorf("не удалось преобразовать duration: %w", err)
	}
	return event, nil
}

func parsePostgresInterval(interval string) (time.Duration, error) {
	// Парсим строку в формате "HH:MM:SS"
	parts := strings.Split(interval, ":")
	if len(parts) != 3 {
		return 0, errors.New("invalid interval format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	duration := time.Duration(hours*60*60+minutes*60+seconds) * time.Second
	return duration, nil
}
