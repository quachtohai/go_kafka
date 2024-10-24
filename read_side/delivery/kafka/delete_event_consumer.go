package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func (s *readerMessageProcessor) ProcessDeletedEvent(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.commitMessage(ctx, r, m)
	fmt.Printf("DELETE MESSAGE FROM KAFKA")
}
