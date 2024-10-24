package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func (s *eventMessageProcessor) ProcessDeleteEvent(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.commitMessage(ctx, r, m)
}
