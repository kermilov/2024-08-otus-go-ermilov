package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	require.NotNil(t, New("INFO"))
	require.Panics(t, func() { New("PANIC") })
}
