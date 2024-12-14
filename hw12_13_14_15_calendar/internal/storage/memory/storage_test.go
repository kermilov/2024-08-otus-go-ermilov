package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

var tested = New()

func TestStorage(t *testing.T) {
	now := time.Now()

	event1, err := tested.Create(context.Background(), storage.Event{ID: "1", Title: "create"})
	require.Nil(t, err)
	require.NotNil(t, event1)
	require.Equal(t, "1", event1.ID)
	require.Equal(t, "create", event1.Title)

	events, err := tested.FindByDay(context.Background(), now)
	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events, 0)

	event1, err = tested.FindByID(context.Background(), "1")
	require.Nil(t, err)
	require.NotNil(t, event1)
	require.Equal(t, "1", event1.ID)
	require.Equal(t, "create", event1.Title)

	err = tested.Update(context.Background(), "1", storage.Event{ID: "1", Title: "update", DateTime: now})
	require.Nil(t, err)

	events, err = tested.FindByDay(context.Background(), now)
	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events, 1)
	require.Equal(t, "1", events[0].ID)
	require.Equal(t, "update", events[0].Title)

	events, err = tested.FindByWeek(context.Background(), now)
	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events, 1)
	require.Equal(t, "1", events[0].ID)
	require.Equal(t, "update", events[0].Title)

	events, err = tested.FindByMonth(context.Background(), now)
	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events, 1)
	require.Equal(t, "1", events[0].ID)
	require.Equal(t, "update", events[0].Title)

	event1, err = tested.FindByID(context.Background(), "1")
	require.Nil(t, err)
	require.NotNil(t, event1)
	require.Equal(t, "1", event1.ID)
	require.Equal(t, "update", event1.Title)

	err = tested.Update(context.Background(), "1", storage.Event{ID: "1", Title: "update", DateTime: now.Add(-storage.Day)})
	require.Nil(t, err)

	events, err = tested.FindByDay(context.Background(), now)
	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events, 0)

	err = tested.Update(context.Background(), "1", storage.Event{ID: "1", Title: "update", DateTime: now.Add(-storage.Week)})
	require.Nil(t, err)

	events, err = tested.FindByWeek(context.Background(), now)
	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events, 0)

	err = tested.Update(context.Background(), "1", storage.Event{ID: "1", Title: "update", DateTime: now.Add(-storage.Week * 4)})
	require.Nil(t, err)

	events, err = tested.FindByMonth(context.Background(), now)
	require.Nil(t, err)
	require.NotNil(t, events)
	require.Len(t, events, 0)

	err = tested.Delete(context.Background(), "1")
	require.Nil(t, err)

	event1, err = tested.FindByID(context.Background(), "1")
	require.NotNil(t, err)
}
