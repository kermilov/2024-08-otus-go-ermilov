package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/app"
	internalСonsumer "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/consumer"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/consumer/kafka"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_storer.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	logg := logger.New(config.Logger.Level)

	storage := getStorage(config)
	calendar := app.New(logg, storage)

	// Создаем HTTP сервер для метрик
	metricsServer := &http.Server{
		Addr:              config.HTTP.String(),
		Handler:           promhttp.Handler(),
		ReadHeaderTimeout: 10 * time.Second, // от G112: Potential Slowloris Attack
	}

	consumer := getConsumer(logg, config, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := consumer.Stop(ctx); err != nil {
			logg.Error("failed to stop consumer: " + err.Error())
		}
	}()

	logg.Info("calendar storer is running...")

	go func() {
		if err := consumer.Start(ctx); err != nil {
			logg.Error("failed to start consumer: " + err.Error())
			cancel()
			return
		}
	}()
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Error("failed to start metrics server: " + err.Error())
		}
	}()
	<-ctx.Done()
}

func getStorage(config Config) app.Storage {
	if _, isOk := supportedStorages[config.Storage]; !isOk {
		panic(fmt.Errorf("неизвестный тип хранения: %s", config.Storage))
	}
	switch config.Storage {
	case InMemoryStorage:
		return memorystorage.New()
	case SQLStorage:
		return sqlstorage.New(config.DB.String())
	}
	return nil
}

func getConsumer(logg *logger.Logger, config Config, calendar *app.App) internalСonsumer.Consumer {
	if _, isOk := supportedMessageBrokers[config.MessageBroker]; !isOk {
		panic(fmt.Errorf("неизвестный брокер сообщений: %s", config.MessageBroker))
	}
	switch config.MessageBroker { //nolint: gocritic
	case Kafka:
		return kafka.NewConsumer(
			logg, calendar, config.Kafka.String(), config.NotificationQueue, internalСonsumer.SaveNotification,
		)
	}
	return nil
}
