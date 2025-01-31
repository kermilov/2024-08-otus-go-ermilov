package util

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
)

const (
	timeOut = 10 * time.Second
)

type Logger interface {
	Warning(msg string)
}

func KafkaCheckConnect(ctx context.Context, addr string, logger Logger, topic string) error {
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
			conn, err = dialer.DialContext(ctx, "tcp", addr)
			if err != nil {
				logger.Warning(fmt.Sprintf("ожидание подключения к Apache Kafka: %s", err.Error()))
				return err
			}
			return nil
		},
		expBackoff,
	)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к Apache Kafka по адресу %s: %w", addr, err)
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
	return nil
}
