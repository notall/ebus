eBus
===============
eBus is a lightweight in-memory pubsub EventBus for GoLang.
Subscribers can execute asynchronously, and can specify a thread pool to perform tasks.

#### Installation
```shell
go get github.com/notall/ebus
```

#### Import package
```go
import "github.com/notall/ebus"
```

#### Example
```go
type event struct{}

func main() {
    // New(WithLogger(logger), WithThreadPool(pool))
    bus := New()
    topic := "topic"
    handler := func(ctx context.Context, arg interface{}) {
        e, ok := arg.(*event)
        // your code here...
    }
    bus.Subscribe(topic, handler);

    ev1 := &event{}
    bus.Publish(context.Background(), topic, ev1);

    ev2 := &event{}
    w := bus.PublishAsync(context.Background(), topic, ev2);
    w.WaitComplete()
}
```

#### Implemented methods
* **New()**
* **Subscribe()**
* **SubscribeOnce()**
* **Unsubscribe()**
* **Publish()**
* **PublishAsync()**

#### New() EventBus
New Returns new EventBus with empty handlers
```go
bus := eBus.New()
```
Specify thread pool and logger
```go
import (
    "github.com/panjf2000/ants/v2"
)

type _logger struct {
}

func (log *_logger) CtxError(ctx context.Context, format string, v ...interface{}) {
    fmt.Printf(format, v...)
}

pool, err := ants.NewPool(10000)
if err != nil {
    panic(err)
}
logger := new(_logger)
bus := eBus.New(eBus.WithLogger(logger), eBus.WithThreadPool(pool))
```

#### Subscribe() error
Subscribe to a topic.
Returns error if duplicated handler subscribed to same topic.
```go
handler := func(ctx context.Context, arg interface{}) {
    // your code here...
}
bus.Subscribe("topic", handler)
```

#### SubscribeOnce() error
Subscribe to a topic once, handler will be removed after executing.
Returns error if duplicated handler subscribed to same topic.
```go
handler := func(ctx context.Context, arg interface{}) {
    // your code here...
}
bus.SubscribeOnce("topic", handler)
```

#### Unsubscribe() error
Remove handler for a topic.
Returns error if there are no handlers subscribed to the topic.
```go
bus.Unsubscribe("topic", handler)
```

#### Publish()
Executes handlers subscribed for a topic.
```go
type event struct {
    val int
}
e := &event{
    val: 5,
}
bus.Publish(context.Background(), "topic", e)
```

#### PublishAsync()
Executes handlers subscribed for a topic.
Returns Waiter can be use to wait for all async handlers to complete
```go
type event struct {
    val int
}
e := &event{
    val: 5,
}
w := bus.PublishAsync(context.Background(), "topic", e)
w.WaitComplete()
```