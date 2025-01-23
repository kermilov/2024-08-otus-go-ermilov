package kafka

import (
	"context"
	"encoding/json"

	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/producer"
	"github.com/kermilov/2024-08-otus-go-ermilov/hw12_13_14_15_calendar/internal/util"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	logger producer.Logger
	addr   string
	topic  string
	writer *kafka.Writer
}

func NewProducer(logger producer.Logger, addr, topic string) *Producer {
	return &Producer{
		logger: logger,
		addr:   addr,
		topic:  topic,
	}
}

func (p *Producer) Start(ctx context.Context) error {
	err := util.KafkaCheckConnect(ctx, p.addr, p.logger, p.topic)
	if err != nil {
		return err
	}
	p.writer = &kafka.Writer{
		Addr:     kafka.TCP(p.addr),
		Topic:    p.topic,
		Balancer: &kafka.LeastBytes{},
	}
	return nil
}

func (p *Producer) Stop(_ context.Context) error {
	return p.writer.Close()
}

func (p *Producer) Send(ctx context.Context, notification *util.Notification) error {
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
	return nil
}
