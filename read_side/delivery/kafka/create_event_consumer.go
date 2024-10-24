package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/segmentio/kafka-go"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}

func (s *readerMessageProcessor) ProcessCreatedEvent(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.commitMessage(ctx, r, m)
	fmt.Printf("CREATED EVENT MESSAGE FROM KAFKA")
}
