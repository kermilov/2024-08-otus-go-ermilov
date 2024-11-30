package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const toPath = "testdata/out.txt"

func TestErrNotExist(t *testing.T) {
	err := Copy("testdata/notexists.txt", toPath, 0, 0)
	require.NotNil(t, err)
	require.True(t, os.IsNotExist(err))
}

func TestErrOffsetExceedsFileSize(t *testing.T) {
	err := Copy("testdata/input.txt", toPath, 1000000000000, 0)
	require.NotNil(t, err)
	require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
}

func TestErrSamePath(t *testing.T) {
	err := Copy("testdata/input.txt", "testdata/input.txt", 1000000000000, 0)
	require.NotNil(t, err)
	require.True(t, errors.Is(err, ErrSamePath))
}

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
