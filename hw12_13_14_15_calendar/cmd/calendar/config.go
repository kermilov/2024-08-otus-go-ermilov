package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	InMemoryStorage = "in-memory"
	SQLStorage      = "sql"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf `json:"logger"`
	Storage string     `json:"storage"`
	DB      DBConf     `json:"db"`
}

type LoggerConf struct {
	Level string `json:"level"`
}

type DBConf struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Schema   string `json:"schema"`
}

func NewConfig() Config {
	file, err := os.ReadFile(configFile)
	if err != nil {
		panic(fmt.Errorf("не удалось прочитать файл конфигурации: %w", err))
	}
	config := Config{}

	if err := json.Unmarshal(file, &config); err != nil {
		panic(fmt.Errorf("не удалось распарсить файл конфигурации: %w", err))
	}
	return config
}
