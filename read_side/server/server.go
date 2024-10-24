package server

import (
	"context"
	kafkaClient "go_kafka/pkg/kafka"
	"go_kafka/pkg/logger"
	"go_kafka/read_side/config"
	readerKafka "go_kafka/read_side/delivery/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type server struct {
	log       logger.Logger
	cfg       *config.KafkaConfig
	kafkaConn *kafka.Conn
}

func NewServer(log logger.Logger, cfg *config.KafkaConfig) *server {

	return &server{log: log, cfg: cfg}
}
func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	readerMessageProcessor := readerKafka.NewReaderMessageProcessor(s.log, s.cfg)
	s.log.Info("Starting Reader Kafka consumers")
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), readerKafka.PoolSize, readerMessageProcessor.ProcessMessages)
	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errcheck

	<-ctx.Done()
	return nil
}
