package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

type EnvValue struct {
	Value      string
	NeedRemove bool
}

func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := filepath.Base(file)
		if strings.Contains(name, "=") {
			return nil, fmt.Errorf("invalid env file name: %s", name)
		}

		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		data = bytes.TrimSpace(data)
		if len(data) == 0 {
			env[name] = EnvValue{NeedRemove: true}
			continue
		}

		data = bytes.ReplaceAll(data, []byte{0x00}, []byte{'\n'})
		value := string(data)

		env[name] = EnvValue{Value: value}
	}

	return env, nil
}
