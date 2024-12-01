package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "envdir_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		files    map[string]string
		expected Environment
		wantErr  bool
	}{
		{
			name: "simple",
			files: map[string]string{
				"FOO": "bar",
				"BAR": "baz",
			},
			expected: Environment{
				"FOO": EnvValue{Value: "bar"},
				"BAR": EnvValue{Value: "baz"},
			},
		},
		{
			name: "ignore second line",
			files: map[string]string{
				"FOO": "bar\nPLEASE IGNORE SECOND LINE",
				"BAR": "baz",
			},
			expected: Environment{
				"FOO": EnvValue{Value: "bar"},
				"BAR": EnvValue{Value: "baz"},
			},
		},
		{
			name: "not trim leading spaces",
			files: map[string]string{
				"FOO": "   foo",
				"BAR": "baz",
			},
			expected: Environment{
				"FOO": EnvValue{Value: "   foo"},
				"BAR": EnvValue{Value: "baz"},
			},
		},
		{
			name: "trim trailing spaces",
			files: map[string]string{
				"FOO": "   foo  ",
				"BAR": "baz",
			},
			expected: Environment{
				"FOO": EnvValue{Value: "   foo"},
				"BAR": EnvValue{Value: "baz"},
			},
		},
		{
			name: "trim trailing tab",
			files: map[string]string{
				"FOO": "   foo\t",
				"BAR": "baz",
			},
			expected: Environment{
				"FOO": EnvValue{Value: "   foo"},
				"BAR": EnvValue{Value: "baz"},
			},
		},
		{
			name: "not ignore line after zero byte",
			files: map[string]string{
				"FOO": "   foo" + string([]byte{0x00}) + "with new line",
				"BAR": "baz",
			},
			expected: Environment{
				"FOO": EnvValue{Value: "   foo\nwith new line"},
				"BAR": EnvValue{Value: "baz"},
			},
		},
		{
			name: "empty file",
			files: map[string]string{
				"FOO": "",
			},
			expected: Environment{
				"FOO": EnvValue{NeedRemove: true},
			},
		},
		{
			name: "file with =",
			files: map[string]string{
				"FOO=BAR": "baz",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0o644)
				require.NoError(t, err)
			}

			env, err := ReadDir(tmpDir)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				for k := range env {
					if _, ok := tt.files[k]; !ok {
						delete(env, k)
					}
				}
				require.Equal(t, tt.expected, env)
			}
		})
	}
}
