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
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/producer"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/producer/kafka"
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

	producer := getProducer(logg, calendar, config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := producer.Stop(ctx); err != nil {
			logg.Error("failed to stop producer: " + err.Error())
		}
	}()

	logg.Info("calendar scheduler is running...")

	go func() {
		if err := producer.Start(ctx); err != nil {
			logg.Error("failed to start producer: " + err.Error())
			cancel()
			return
		}
		// Немедленное выполнение задачи при запуске
		producer.ScheduledProcess(ctx)
		// Запуск бесконечного цикла для выполнения задачи каждые 15 минут
		for range time.NewTicker(15 * time.Minute).C {
			producer.ScheduledProcess(ctx)
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

func getProducer(logg *logger.Logger, calendar *app.App, config Config) producer.Producer {
	kafkaProducer := kafka.NewProducer(logg, calendar, config.Kafka.String())
	return kafkaProducer
}
