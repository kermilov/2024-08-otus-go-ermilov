package scheduler

import (
	"context"
	"time"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/util"
	"github.com/prometheus/client_golang/prometheus"
)

// Общий интерфейс логгера на разные реализации планировщика.
type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Producer interface {
	Start(ctx context.Context) error
	Send(ctx context.Context, notification *util.Notification) error
	Stop(ctx context.Context) error
}

type Application interface {
	FindForSendNotification(ctx context.Context, date time.Time) ([]storage.Event, error)
	SetIsSendNotification(ctx context.Context, ids []string) error
	DeleteOldEvents(ctx context.Context, date time.Time) error
}

type Scheduler struct {
	logger   Logger
	app      Application
	producer Producer
	ticker   *time.Ticker
	duration time.Duration
}

// Определяем свои метрики.
var sendNotificationsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "send_notifications_total",
		Help: "Total number of send notifications",
	},
)

func init() {
	// Регистрируем метрики
	prometheus.MustRegister(sendNotificationsTotal)
}

func NewScheduler(logger Logger, app Application, producer Producer, duration time.Duration) *Scheduler {
	return &Scheduler{logger: logger, app: app, producer: producer, duration: duration}
}

func (p *Scheduler) Start(ctx context.Context) error {
	if err := p.producer.Start(ctx); err != nil {
		return err
	}
	// Немедленное выполнение задачи при запуске
	err := p.scheduledProcess(ctx)
	if err != nil {
		return err
	}
	// Запуск бесконечного цикла для выполнения задачи
	p.ticker = time.NewTicker(p.duration)
	for range p.ticker.C {
		p.scheduledProcess(ctx)
	}
	return nil
}

func (p *Scheduler) Stop(ctx context.Context) error {
	p.ticker.Stop()
	if err := p.producer.Stop(ctx); err != nil {
		p.logger.Error("failed to stop producer: " + err.Error())
		return err
	}
	return nil
}

func (p *Scheduler) scheduledProcess(ctx context.Context) error {
	p.logger.Info("send notifications and delete old events - start")
	if err := p.sendNotifications(ctx); err != nil {
		p.logger.Error("failed to send notifications: " + err.Error())
		return err
	}
	if err := p.deleteOldEvents(ctx); err != nil {
		p.logger.Error("failed to delete old events: " + err.Error())
		return err
	}
	p.logger.Info("send notifications and delete old events - end")
	return nil
}

func (p *Scheduler) sendNotifications(ctx context.Context) error {
	events, err := p.app.FindForSendNotification(ctx, time.Now())
	if err != nil {
		return err
	}
	for _, v := range events {
		notification := &util.Notification{
			ID:       v.ID,
			Title:    v.Title,
			DateTime: v.DateTime,
			UserID:   v.UserID,
		}
		err := p.producer.Send(ctx, notification)
		if err != nil {
			return err
		}
		p.app.SetIsSendNotification(ctx, []string{v.ID})
		// Инкрементируем счетчик посланных уведомлений
		sendNotificationsTotal.Inc()
	}
	return nil
}

func (p *Scheduler) deleteOldEvents(ctx context.Context) error {
	return p.app.DeleteOldEvents(ctx, time.Now())
}
