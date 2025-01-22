package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	InMemoryStorage = "in-memory"
	SQLStorage      = "sql"
)

var supportedStorages = map[string]struct{}{
	InMemoryStorage: {},
	SQLStorage:      {},
}

const (
	Kafka = "kafka"
)

var supportedMessageBrokers = map[string]struct{}{
	Kafka: {},
}

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger            LoggerConf     `json:"logger"`
	Storage           string         `json:"storage"`
	NotificationQueue string         `json:"notificationQueue"`
	MessageBroker     string         `json:"messageBroker"`
	DB                DBConf         `json:"db"`
	Kafka             KafkaConf      `json:"kafka"`
	Duration          CustomDuration `json:"duration"`
}

type CustomDuration struct {
	time.Duration
}

func (d *CustomDuration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
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

type KafkaConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (c *KafkaConf) String() string {
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
