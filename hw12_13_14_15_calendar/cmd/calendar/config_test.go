package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	configFile = "../../configs/config.json"
	actual := NewConfig()
	require.Equal(t, "INFO", actual.Logger.Level)
}
