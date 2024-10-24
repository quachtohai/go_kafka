package kafka

import (
	"context"
	"fmt"
	"go_kafka/pkg/logger"
	"go_kafka/read_side/config"

	"sync"

	"github.com/segmentio/kafka-go"
)

const (
	PoolSize = 30
)

type readerMessageProcessor struct {
	log logger.Logger
	cfg *config.KafkaConfig
}

func NewReaderMessageProcessor(log logger.Logger, cfg *config.KafkaConfig) *readerMessageProcessor {
	return &readerMessageProcessor{log: log, cfg: cfg}
}

func (s *readerMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
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
		case s.cfg.KafkaTopics.EventCreated.TopicName:
			s.ProcessCreatedEvent(ctx, r, m)
		case s.cfg.KafkaTopics.EventUpdated.TopicName:
			s.ProcessUpdatedEvent(ctx, r, m)
		case s.cfg.KafkaTopics.EventDeleted.TopicName:
			s.ProcessDeletedEvent(ctx, r, m)
		}
	}
}
