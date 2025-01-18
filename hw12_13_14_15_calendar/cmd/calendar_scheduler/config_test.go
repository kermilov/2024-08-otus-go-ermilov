package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	configFile = "../../configs/calendar_scheduler.json"
	actual := NewConfig()
	require.Equal(t, "INFO", actual.Logger.Level)
	require.Equal(t, SQLStorage, actual.Storage)
	require.Equal(t, "localhost", actual.DB.Host)
	require.Equal(t, 5432, actual.DB.Port)
	require.Equal(t, "postgres", actual.DB.User)
	require.Equal(t, "postgres", actual.DB.Password)
	require.Equal(t, "otus", actual.DB.Name)
	require.Equal(t, "calendar", actual.DB.Schema)
	require.Equal(t, "localhost", actual.Kafka.Host)
	require.Equal(t, 29092, actual.Kafka.Port)
}
