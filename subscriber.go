package ebus

import "context"

// Subscriber subscribe topic with handler
type Subscriber interface {
	// Returns error if duplicated handler subscribed to same topic.
	Subscribe(topic string, handler Handler) error
	// SubscribeOnce handler will be removed after executing
	// Returns error if duplicated handler subscribed to same topic.
	SubscribeOnce(topic string, handler Handler) error
	// Unsubscribe removes handler for a topic
	// Returns error if there are no handlers subscribed to the topic.
	Unsubscribe(topic string, handler Handler) error
}

// Handler function to handle event
type Handler func(ctx context.Context, event interface{})
