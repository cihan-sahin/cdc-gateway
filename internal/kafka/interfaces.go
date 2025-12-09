package kafka

import "context"

type Consumer interface {
	Consume(ctx context.Context, handler func(ctx context.Context, msg Message) error) error
	Close() error
}

type Producer interface {
	Send(ctx context.Context, msg Message) error
	Close() error
}
