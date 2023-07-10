package ebus

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// New returns new EventBus with empty handlers.
func New(options ...Option) EventBus {
	bus := &_eventBus{
		handlers: make(map[string][]*eventHandler),
		lock:     sync.Mutex{},
	}
	opts := loadOptions(options...)
	if opts.Logger == nil {
		opts.Logger = defaultLogger
	}
	if opts.ThreadPool == nil {
		opts.ThreadPool = defaultThreadPool
	}
	bus.opts = opts
	return bus
}

// EventBus pubsub
type EventBus interface {
	Subscriber
	Publisher
}

type _eventBus struct {
	handlers map[string][]*eventHandler
	lock     sync.Mutex
	opts     *_options
}

type eventHandler struct {
	handler  reflect.Value
	flagOnce bool
}

func (bus *_eventBus) Subscribe(topic string, handler Handler) error {
	return bus.doSubscribe(topic, handler, false)
}

func (bus *_eventBus) SubscribeOnce(topic string, handler Handler) error {
	return bus.doSubscribe(topic, handler, true)
}

func (bus *_eventBus) Unsubscribe(topic string, handler Handler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if len(bus.handlers[topic]) > 0 {
		bus.removeHandler(topic, bus.findHandlerIdx(topic, reflect.ValueOf(handler)))
		return nil
	}
	return ErrorTopicNotExist
}

func (bus *_eventBus) Publish(ctx context.Context, topic string, event interface{}) {
	bus.doPublish(ctx, topic, event, false)
}

func (bus *_eventBus) PublishAsync(ctx context.Context, topic string, event interface{}) Waiter {
	wg := bus.doPublish(ctx, topic, event, true)
	return &_waiter{
		wg: wg,
	}
}

func (bus *_eventBus) doSubscribe(topic string, handler Handler, flagOnce bool) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	idx := bus.findHandlerIdx(topic, reflect.ValueOf(handler))
	if idx > -1 {
		return ErrorDuplicatedHandler
	}
	eventHandler := &eventHandler{
		handler:  reflect.ValueOf(handler),
		flagOnce: flagOnce,
	}
	bus.handlers[topic] = append(bus.handlers[topic], eventHandler)
	return nil
}

func (bus *_eventBus) doPublish(ctx context.Context, topic string, event interface{}, async bool) *sync.WaitGroup {
	handlers := bus.getTopicHandlersAndRemoveOnce(topic)
	if len(handlers) <= 0 {
		return nil
	}

	var waitGroup *sync.WaitGroup
	if async {
		waitGroup = &sync.WaitGroup{}
	}
	for _, v := range handlers {
		if !async {
			bus.handleEvent(ctx, v.handler, event)
			continue
		}

		waitGroup.Add(1)
		vHandler := v.handler
		err := bus.opts.ThreadPool.Submit(
			func() {
				defer waitGroup.Done()
				bus.handleEvent(ctx, vHandler, event)
			},
		)
		// If submitting a task to the pool fails, it will be executed in the current thread
		if err != nil {
			bus.opts.Logger.CtxError(ctx, "[ebus] submit task failed, topic:%v error:%v", topic, err)
			bus.handleEvent(ctx, vHandler, event)
			waitGroup.Done()
		}
	}
	return waitGroup
}

func (bus *_eventBus) handleEvent(ctx context.Context, handler reflect.Value, event interface{}) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			bus.opts.Logger.CtxError(ctx, "[ebus] panic, error:%v stack:%v", err, "...\n"+string(buf))
		}
	}()

	var v reflect.Value
	if event != nil {
		v = reflect.ValueOf(event)
	} else {
		v = reflect.New(handler.Type().In(1)).Elem()
	}
	handler.Call([]reflect.Value{reflect.ValueOf(ctx), v})
}

func (bus *_eventBus) getTopicHandlersAndRemoveOnce(topic string) []*eventHandler {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	handlers := bus.handlers[topic]
	if len(handlers) <= 0 {
		return nil
	}
	copyHandlers := make([]*eventHandler, len(handlers))
	copy(copyHandlers, handlers)
	for _, v := range copyHandlers {
		if v.flagOnce {
			bus.removeHandler(topic, bus.findHandlerIdx(topic, v.handler))
		}
	}
	return copyHandlers
}

func (bus *_eventBus) removeHandler(topic string, idx int) {
	l := len(bus.handlers[topic])
	if idx >= l || idx < 0 {
		return
	}
	copy(bus.handlers[topic][idx:], bus.handlers[topic][idx+1:])
	bus.handlers[topic][l-1] = nil
	bus.handlers[topic] = bus.handlers[topic][:l-1]
}

func (bus *_eventBus) findHandlerIdx(topic string, handler reflect.Value) int {
	if len(bus.handlers[topic]) <= 0 {
		return -1
	}
	for idx, v := range bus.handlers[topic] {
		if v.handler.Pointer() == handler.Pointer() {
			return idx
		}
	}
	return -1
}

type _waiter struct {
	wg *sync.WaitGroup
}

func (w *_waiter) WaitComplete() {
	w.wg.Wait()
}
