package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const toPath = "testdata/out.txt"

func TestCopy(t *testing.T) {
	tests := []struct {
		offset           int64
		limit            int64
		expectedFilePath string
	}{
		{0, 0, "testdata/out_offset0_limit0.txt"},
		{0, 10, "testdata/out_offset0_limit10.txt"},
		{0, 1000, "testdata/out_offset0_limit1000.txt"},
		{0, 10000, "testdata/out_offset0_limit10000.txt"},
		{100, 1000, "testdata/out_offset100_limit1000.txt"},
		{6000, 1000, "testdata/out_offset6000_limit1000.txt"},
	}
	for _, tc := range tests {
		err := Copy("testdata/input.txt", toPath, tc.offset, tc.limit)
		require.Nil(t, err)
		actualContent, err := os.ReadFile(toPath)
		require.Nil(t, err)
		expectedContent, err := os.ReadFile(tc.expectedFilePath)
		require.Nil(t, err)
		require.Equal(t, expectedContent, actualContent)
		err = os.Remove(toPath)
		require.Nil(t, err)
	}
}
