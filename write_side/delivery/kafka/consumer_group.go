package kafka

import (
	"context"
	"fmt"
	"go_kafka/pkg/logger"
	"go_kafka/write_side/config"

	"sync"

	"github.com/segmentio/kafka-go"
)

const (
	PoolSize = 30
)

type eventMessageProcessor struct {
	log logger.Logger
	cfg *config.KafkaConfig
}

func NewEventMessageProcessor(log logger.Logger, cfg *config.KafkaConfig) *eventMessageProcessor {
	return &eventMessageProcessor{log: log, cfg: cfg}
}

func (s *eventMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)

		if err != nil {
			s.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}
		fmt.Printf(m.Topic)
		s.logProcessMessage(m, workerID)

		switch m.Topic {
		case s.cfg.KafkaTopics.EventCreate.TopicName:
			s.ProcessCreateEvent(ctx, r, m)
		case s.cfg.KafkaTopics.EventUpdate.TopicName:
			s.ProcessUpdateEvent(ctx, r, m)
		case s.cfg.KafkaTopics.EventDelete.TopicName:
			s.ProcessDeleteEvent(ctx, r, m)
		}
	}
}
