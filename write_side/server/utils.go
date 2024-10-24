package server

import (
	"context"
	kafkaClient "go_kafka/pkg/kafka"
	"net"
	"strconv"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

const (
	stackSize = 1 << 10 // 1 KB
)

func (s *server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	s.kafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.log.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (s *server) initKafkaTopics(ctx context.Context) {
	controller, err := s.kafkaConn.Controller()
	if err != nil {
		s.log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	s.log.Infof("established new kafka controller connection: %s", controllerURI)

	eventCreateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.EventCreate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.EventCreate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.EventCreate.ReplicationFactor,
	}

	eventCreatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.EventCreated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.EventCreated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.EventCreated.ReplicationFactor,
	}

	eventUpdateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.EventUpdate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.EventUpdate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.EventUpdate.ReplicationFactor,
	}

	eventUpdatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.EventUpdated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.EventUpdated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.EventUpdated.ReplicationFactor,
	}

	eventDeleteTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.EventDelete.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.EventDelete.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.EventDelete.ReplicationFactor,
	}

	eventDeletedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.EventDeleted.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.EventDeleted.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.EventDeleted.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		eventCreateTopic,
		eventUpdateTopic,
		eventCreatedTopic,
		eventUpdatedTopic,
		eventDeleteTopic,
		eventDeletedTopic,
	); err != nil {
		s.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	s.log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{eventCreateTopic, eventUpdateTopic,
		eventCreatedTopic, eventUpdatedTopic, eventDeleteTopic, eventDeletedTopic})
}

func (s *server) getConsumerGroupTopics() []string {
	return []string{
		s.cfg.KafkaTopics.EventCreate.TopicName,
		s.cfg.KafkaTopics.EventUpdate.TopicName,
		s.cfg.KafkaTopics.EventDelete.TopicName,
	}
}
