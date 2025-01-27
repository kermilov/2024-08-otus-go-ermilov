package integrationtests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

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

	// Удаление события
	req, err = http.NewRequestWithContext(ctx, "DELETE", "http://calendar-app:8080/event/"+eventID, nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
