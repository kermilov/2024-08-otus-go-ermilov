package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/viper"
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
	HTTP              HTTPServerConf `json:"http"`
	Kafka             KafkaConf      `json:"kafka"`
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

type KafkaConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (c *KafkaConf) String() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func NewConfig() Config {
	// Указываем полный путь к файлу конфигурации
	viper.SetConfigFile(configFile)

	// Чтение переменных окружения
	viper.AutomaticEnv() // Автоматически связывает переменные окружения с конфигурацией

	// Связки для logger
	viper.BindEnv("logger.level", "LOGGER_LEVEL")

	// Связки для storage
	viper.BindEnv("storage", "STORAGE_TYPE")

	// Связки для notificationQueue
	viper.BindEnv("notificationQueue", "NOTIFICATION_QUEUE")

	// Связки для messageBroker
	viper.BindEnv("messageBroker", "MESSAGE_BROKER")

	// Связки для db
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.name", "DB_NAME")
	viper.BindEnv("db.schema", "DB_SCHEMA")

	// Связки для http
	viper.BindEnv("http.host", "HTTP_HOST")
	viper.BindEnv("http.port", "HTTP_PORT")

	// Связки для kafka
	viper.BindEnv("kafka.host", "KAFKA_HOST")
	viper.BindEnv("kafka.port", "KAFKA_PORT")

	// Чтение конфигурации из файла
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("не удалось прочитать файл конфигурации: %w", err))
	}

	// Загрузка конфигурации в структуру
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("не удалось распарсить файл конфигурации: %w", err))
	}
	return config
}
