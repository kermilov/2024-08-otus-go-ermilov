package kafka

import (
	"context"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/consumer"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/util"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	logger   consumer.Logger
	app      consumer.Application
	addr     string
	topic    string
	reader   *kafka.Reader
	listener consumer.ListenerFunc
}

func NewConsumer(
	logger consumer.Logger, app consumer.Application, addr, topic string, listener consumer.ListenerFunc,
) *Consumer {
	return &Consumer{
		logger:   logger,
		app:      app,
		addr:     addr,
		topic:    topic,
		listener: listener,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	err := util.KafkaCheckConnect(ctx, c.addr, c.logger, c.topic)
	if err != nil {
		return err
	}
	// Создание нового потребителя
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.addr},
		Topic:   c.topic,
	})
	// Чтение сообщений
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.logger.Error(err.Error())
			continue
		}
		err = c.listener(ctx, c.app, msg.Value)
		if err != nil {
			c.logger.Error(err.Error())
			continue
		}
	}
}

func (c *Consumer) Stop(_ context.Context) error {
	c.reader.Close()
	return nil
}
