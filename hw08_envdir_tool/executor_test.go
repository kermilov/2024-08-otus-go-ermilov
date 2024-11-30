package main

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name     string
		cmd      []string
		env      Environment
		expected int
	}{
		{
			name:     "simple",
			cmd:      []string{"echo", "hello"},
			env:      Environment{},
			expected: 0,
		},
		{
			name:     "with env",
			cmd:      []string{"echo", "$FOO"},
			env:      Environment{"FOO": EnvValue{Value: "bar"}},
			expected: 0,
		},
		{
			name:     "with error",
			cmd:      []string{"non-existent-command"},
			env:      Environment{},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnCode := RunCmd(tt.cmd, tt.env)
			if returnCode != tt.expected {
				t.Errorf("RunCmd() = %d, want %d", returnCode, tt.expected)
			}
		})
	}
}
