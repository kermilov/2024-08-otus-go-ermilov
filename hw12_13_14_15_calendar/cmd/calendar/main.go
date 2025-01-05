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
	internalgrpc "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "C:/Users/kermilov/Otus/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/configs/config.json", "Path to configuration file")
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

	serverHttp := internalhttp.NewServer(logg, calendar, config.HTTP.String())
	serverGrpc := internalgrpc.NewServer(logg, calendar, config.GRPC.String())

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHttp.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverGrpc.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := serverHttp.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()
	go func() {
		if err := serverGrpc.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
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
