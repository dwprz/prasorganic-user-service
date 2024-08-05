package client

import (
	"context"
)

type RabbitMQ interface {
	Publish(ctx context.Context, exchange string, key string, message any)
	Close()
}
