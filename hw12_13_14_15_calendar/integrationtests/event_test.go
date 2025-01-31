package integrationtests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/logger"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/util"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func TestEventLifecycle(t *testing.T) {
	time.Sleep(10 * time.Second) // Задержка для инициализации сервисов

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Добавление события
	const eventID = "10c547d6-c702-4dfc-a6e1-42d54c913825"
	event := map[string]interface{}{
		"id":                   eventID,
		"title":                "Test Event",
		"datetime":             "2025-01-01T10:00:00Z",
		"duration":             "2h",
		"userid":               1,
		"notificationDuration": "30m",
	}
	eventBytes, _ := json.Marshal(event)
	req, err := http.NewRequestWithContext(ctx, "POST", "http://calendar-app:8080/event", bytes.NewBuffer(eventBytes))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверка бизнес ошибки на добавление события в занятое время
	req, err = http.NewRequestWithContext(ctx, "POST", "http://calendar-app:8080/event", bytes.NewBuffer(eventBytes))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	// Получение листинга событий на день
	req, err = http.NewRequestWithContext(ctx, "GET", "http://calendar-app:8080/events?day=2025-01-01T00:00:00Z", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Получение листинга событий на неделю
	req, err = http.NewRequestWithContext(ctx, "GET", "http://calendar-app:8080/events?week=2024-12-30T00:00:00Z", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Получение листинга событий на месяц
	req, err = http.NewRequestWithContext(ctx, "GET", "http://calendar-app:8080/events?month=2025-01-01T00:00:00Z", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Отправка уведомлений
	// Чтение уведомления из Kafka
	err = util.KafkaCheckConnect(ctx, "kafka:9092", logger.New("WARNING"), "notification-topic")
	require.NoError(t, err)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "notification-topic",
	})

	defer r.Close()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message, err := r.ReadMessage(ctx)
	require.NoError(t, err)
	require.Contains(t, string(message.Value), eventID)

	// Проверка наличия записи в таблице notification
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/otus?sslmode=disable&search_path=calendar")
	require.NoError(t, err)
	defer db.Close()

	var count int
	err = db.QueryRow("select count(*) from notification where id = $1", eventID).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)

	// Удаление события
	req, err = http.NewRequestWithContext(ctx, "DELETE", "http://calendar-app:8080/event/"+eventID, nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
