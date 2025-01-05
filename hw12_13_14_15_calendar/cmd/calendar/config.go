package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	InMemoryStorage = "in-memory"
	SQLStorage      = "sql"
)

var supportedStorages = map[string]struct{}{
	InMemoryStorage: {},
	SQLStorage:      {},
}

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf     `json:"logger"`
	Storage string         `json:"storage"`
	DB      DBConf         `json:"db"`
	HTTP    HTTPServerConf `json:"http"`
	GRPC    GRPCServerConf `json:"grpc"`
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

func (db *DBConf) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.Schema)
}

type HTTPServerConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (c *HTTPServerConf) String() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type GRPCServerConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (c *GRPCServerConf) String() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
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
