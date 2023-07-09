package ebus

import "context"

// Publisher publish event
type Publisher interface {
	// Publish executes handlers subscribed for a topic
	Publish(ctx context.Context, topic string, event interface{})
	// PublishAsync asynchronously executes handlers subscribed for a topic
	// Returns Waiter can be use to wait for all async handlers to complete
	PublishAsync(ctx context.Context, topic string, event interface{}) Waiter
}

type Waiter interface {
	// WaitComplete wait for all async handlers to complete
	WaitComplete()
}
