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

func (env Environment) ToEnv() []string {
	envList := make([]string, 0, len(env))
	for name, value := range env {
		if value.NeedRemove {
			continue
		}
		envList = append(envList, name+"="+value.Value)
	}
	return envList
}

func NewEnvironment() Environment {
	env := make(Environment)

	for _, e := range os.Environ() {
		keyValue := strings.Split(e, "=")
		key := keyValue[0]
		value := keyValue[1]

		env[key] = EnvValue{Value: value, NeedRemove: false}
	}
	return env
}

func ReadDir(dir string) (Environment, error) {
	env := NewEnvironment()

	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name := filepath.Base(file)
		if strings.Contains(name, "=") {
			return nil, fmt.Errorf("invalid env file name: %s", name)
		}

		value, err := readFirstLine(file)
		if err != nil {
			return nil, err
		}
		if len(value) == 0 {
			env[name] = EnvValue{NeedRemove: true}
			continue
		}
		env[name] = EnvValue{Value: value}
	}

	return env, nil
}

func readFirstLine(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	data = bytes.Split(data, []byte{'\n'})[0]
	data = bytes.ReplaceAll(data, []byte{0x00}, []byte{'\n'})

	value := string(data)
	return strings.TrimRight(value, " "), nil
}
