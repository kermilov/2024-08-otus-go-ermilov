package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/app"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/logger"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/producer/kafka"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/scheduler"
	memorystorage "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_scheduler.json", "Path to configuration file")
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

	producer := getProducer(logg, config)
	scheduler := scheduler.NewScheduler(logg, calendar, producer, config.Duration.Duration)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := scheduler.Stop(ctx); err != nil {
			logg.Error("failed to stop scheduler: " + err.Error())
		}
	}()

	logg.Info("calendar scheduler is running...")

	go func() {
		if err := scheduler.Start(ctx); err != nil {
			logg.Error("failed to start scheduler: " + err.Error())
			cancel()
			return
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

func getProducer(logg *logger.Logger, config Config) scheduler.Producer {
	if _, isOk := supportedMessageBrokers[config.MessageBroker]; !isOk {
		panic(fmt.Errorf("неизвестный брокер сообщений: %s", config.MessageBroker))
	}
	switch config.MessageBroker { //nolint: gocritic
	case Kafka:
		return kafka.NewProducer(logg, config.Kafka.String(), config.NotificationQueue)
	}
	return nil
}
