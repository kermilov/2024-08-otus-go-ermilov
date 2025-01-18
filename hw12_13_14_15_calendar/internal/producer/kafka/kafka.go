package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/producer"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	logger producer.Logger
	app    producer.Application
	addr   string
	writer *kafka.Writer
}

func NewProducer(logger producer.Logger, app producer.Application, addr string) *Producer {
	return &Producer{
		logger: logger,
		app:    app,
		addr:   addr,
	}
}

const (
	timeOut = 10 * time.Second
	topic   = "event-notification-topic"
)

func (p *Producer) Start(ctx context.Context) error {
	dialer := &kafka.Dialer{
		Timeout: timeOut,
	}
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxInterval = timeOut
	expBackoff.MaxElapsedTime = 0 // Продолжать попытки бесконечно
	var conn *kafka.Conn
	var err error
	err = backoff.Retry(
		func() error {
			conn, err = dialer.DialContext(ctx, "tcp", p.addr)
			if err != nil {
				p.logger.Warning(fmt.Sprintf("ожидание подключения к Apache Kafka: %s", err.Error()))
				return err
			}
			return nil
		},
		expBackoff,
	)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к Apache Kafka по адресу %s: %w", p.addr, err)
	}
	defer conn.Close()
	// Создание топика
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return err
	}
	p.writer = &kafka.Writer{
		Addr:     kafka.TCP(p.addr),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return nil
}

func (p *Producer) Stop(_ context.Context) error {
	return p.writer.Close()
}

func (p *Producer) ScheduledProcess(ctx context.Context) {
	p.logger.Info("send notifications and delete old events - start")
	if err := p.sendNotifications(ctx); err != nil {
		p.logger.Error("failed to send notifications: " + err.Error())
	}
	if err := p.deleteOldEvents(ctx); err != nil {
		p.logger.Error("failed to delete old events: " + err.Error())
	}
	p.logger.Info("send notifications and delete old events - end")
}

func (p *Producer) sendNotifications(ctx context.Context) error {
	events, err := p.app.FindForSendNotification(ctx, time.Now())
	if err != nil {
		return err
	}
	for _, v := range events {
		notification := &producer.Notification{
			ID:       v.ID,
			Title:    v.Title,
			DateTime: v.DateTime,
			UserID:   v.UserID,
		}
		bytes, err := json.Marshal(notification)
		if err != nil {
			return err
		}
		// Сообщение для отправки
		msg := kafka.Message{
			Value: bytes,
		}
		// Отправка сообщения
		err = p.writer.WriteMessages(ctx, msg)
		if err != nil {
			p.logger.Error(err.Error())
			return err
		}
		p.app.SetIsSendNotification(ctx, []string{v.ID})
	}
	return nil
}

func (p *Producer) deleteOldEvents(ctx context.Context) error {
	return p.app.DeleteOldEvents(ctx, time.Now())
}
