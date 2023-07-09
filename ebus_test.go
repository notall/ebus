package ebus

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bus := New()
	if bus == nil {
		t.Log("eventBus not created")
		t.Fail()
	}
}

func TestManySubscribe(t *testing.T) {
	bus := New()
	topic := "topic"
	flag := 0
	fn := func(ctx context.Context, arg interface{}) {
		flag += 1
	}
	err := bus.Subscribe(topic, fn)
	if err != nil {
		t.Log("subscribe failed")
		t.Fail()
	}
	err = bus.Subscribe(topic, fn)
	if err != ErrorDuplicatedHandler {
		t.Log("subscribe failed")
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 1 {
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 2 {
		t.Fail()
	}
}

func TestManySubscribeOnce(t *testing.T) {
	bus := New()
	topic := "topic"
	flag := 0
	fn := func(ctx context.Context, arg interface{}) {
		flag += 1
	}
	err := bus.SubscribeOnce(topic, fn)
	if err != nil {
		t.Log("subscribe failed")
		t.Fail()
	}
	err = bus.SubscribeOnce(topic, fn)
	if err != ErrorDuplicatedHandler {
		t.Log("subscribe failed")
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 1 {
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 1 {
		t.Fail()
	}
}

func TestManySubscribeAndSubscribeOnce(t *testing.T) {
	bus := New()
	topic := "topic"
	flag := 0
	fn := func(ctx context.Context, arg interface{}) {
		flag += 1
	}
	err := bus.Subscribe(topic, fn)
	if err != nil {
		t.Log("subscribe failed")
		t.Fail()
	}
	err = bus.SubscribeOnce(topic, fn)
	if err != ErrorDuplicatedHandler {
		t.Log("subscribe failed")
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 1 {
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 2 {
		t.Fail()
	}
}

func TestManySubscribeOnceAndSubscribe(t *testing.T) {
	bus := New()
	topic := "topic"
	flag := 0
	fn := func(ctx context.Context, arg interface{}) {
		flag += 1
	}
	err := bus.SubscribeOnce(topic, fn)
	if err != nil {
		t.Log("subscribe failed")
		t.Fail()
	}
	err = bus.Subscribe(topic, fn)
	if err != ErrorDuplicatedHandler {
		t.Log("subscribe failed")
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 1 {
		t.Fail()
	}
	bus.Publish(context.Background(), topic, nil)
	if flag != 1 {
		t.Fail()
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := New()
	topic := "topic"
	flag := 0
	fn := func(ctx context.Context, arg interface{}) {
		flag += 1
	}
	bus.Subscribe(topic, fn)
	if bus.Unsubscribe(topic, fn) != nil {
		t.Fail()
	}
	if bus.Unsubscribe(topic, fn) != ErrorTopicNotExist {
		t.Fail()
	}

	bus.SubscribeOnce(topic, fn)
	if bus.Unsubscribe(topic, fn) != nil {
		t.Fail()
	}
	if bus.Unsubscribe(topic, fn) != ErrorTopicNotExist {
		t.Fail()
	}
}

func TestPublish(t *testing.T) {
	type event struct {
		val int
	}
	e := &event{
		val: 5,
	}

	bus := New()
	topic := "topic"
	fn := func(ctx context.Context, arg interface{}) {
		tmp, _ := arg.(*event)
		if tmp != e && e.val != 5 {
			t.Fail()
		}
	}
	bus.Subscribe(topic, fn)
	bus.Publish(context.Background(), "topic", e)
}

func TestPublishPanic(t *testing.T) {
	bus := New()
	topic := "topic"
	result := 0
	fn := func(ctx context.Context, arg interface{}) {
		result++
		panic("panic")
	}
	bus.Subscribe(topic, fn)
	bus.Publish(context.Background(), "topic", 10)
	bus.Publish(context.Background(), "topic", 10)
	if result != 2 {
		t.Fail()
	}
}

func TestPublishWithManySubscription(t *testing.T) {
	bus := New()
	topic := "topic"
	result := 0
	fn1 := func(ctx context.Context, arg interface{}) {
		a := arg.(int)
		result += a
	}

	fn2 := func(ctx context.Context, arg interface{}) {
		a := arg.(int)
		result += a
	}
	bus.Subscribe(topic, fn1)
	bus.Subscribe(topic, fn2)
	bus.Publish(context.Background(), "topic", 10)
	if result != 20 {
		t.Fail()
	}
}

func TestPublishAsync(t *testing.T) {
	bus := New()
	topic := "topic"
	result := 0
	fn1 := func(ctx context.Context, arg interface{}) {
		time.Sleep(time.Second * 2)
		a := arg.(int)
		result = a
	}
	bus.Subscribe(topic, fn1)
	w := bus.PublishAsync(context.Background(), topic, 2)
	w.WaitComplete()
	if result != 2 {
		t.Fail()
	}
}

func TestPublishAsyncManySubscription(t *testing.T) {
	bus := New()
	topic := "topic"
	results := make(chan interface{}, 2)
	fn1 := func(ctx context.Context, arg interface{}) {
		time.Sleep(time.Second * 5)
		results <- "5s"
	}
	fn2 := func(ctx context.Context, arg interface{}) {
		time.Sleep(time.Second * 4)
		results <- "4s"
	}
	bus.Subscribe(topic, fn1)
	bus.Subscribe(topic, fn2)

	w := bus.PublishAsync(context.Background(), topic, nil)

	w.WaitComplete()

	r2 := <-results
	if r2 != "4s" {
		t.Fail()
	}
	r3 := <-results
	if r3 != "5s" {
		t.Fail()
	}
}

func TestPublishAsyncPanic1(t *testing.T) {
	bus := New()
	topic := "topic"
	results := make(chan interface{}, 2)
	fn1 := func(ctx context.Context, arg interface{}) {
		time.Sleep(time.Second * 2)
		results <- "2s"
	}
	fn2 := func(ctx context.Context, arg interface{}) {
		panic("fn2")
	}
	bus.Subscribe(topic, fn1)
	bus.Subscribe(topic, fn2)

	w := bus.PublishAsync(context.Background(), topic, nil)

	time.Sleep(time.Second)
	results <- "1s"

	w.WaitComplete()

	r2 := <-results
	if r2 != "1s" {
		t.Fail()
	}
	r3 := <-results
	if r3 != "2s" {
		t.Fail()
	}
}

func TestPublishAsyncPanic2(t *testing.T) {
	bus := New()
	topic := "topic"
	results := make(chan interface{}, 2)
	fn1 := func(ctx context.Context, arg interface{}) {
		time.Sleep(time.Second * 2)
		results <- "2s"
	}
	fn2 := func(ctx context.Context, arg interface{}) {
		panic("fn2")
	}
	bus.Subscribe(topic, fn1)
	bus.Subscribe(topic, fn2)

	w := bus.PublishAsync(context.Background(), topic, nil)

	time.Sleep(time.Second * 3)
	results <- "3s"

	w.WaitComplete()

	r2 := <-results
	if r2 != "2s" {
		t.Fail()
	}
	r3 := <-results
	if r3 != "3s" {
		t.Fail()
	}
}

func TestPublishAsyncPanic3(t *testing.T) {
	bus := New()
	topic := "topic"
	results := make(chan interface{}, 2)
	fn1 := func(ctx context.Context, arg interface{}) {
		time.Sleep(time.Second * 2)
		results <- "2s"
	}
	fn2 := func(ctx context.Context, arg interface{}) {
		panic("fn2")
	}
	bus.Subscribe(topic, fn1)
	bus.Subscribe(topic, fn2)

	w := bus.PublishAsync(context.Background(), topic, nil)

	w.WaitComplete()
	results <- "0s"

	r2 := <-results
	if r2 != "2s" {
		t.Fail()
	}
	r3 := <-results
	if r3 != "0s" {
		t.Fail()
	}
}
