<img height="250" src="https://raw.githubusercontent.com/qdm12/goshutdown/main/title.svg?sanitize=true">

`goshutdown` is a library to gracefully shutdown your goroutines in your Go program.

## Context

Since Go program are often running multiple goroutines, it is important to shut them down gracefully when the program exits.
Not doing so can result in a loss of data for example, and it's generally a good practice to carefully manage the lifecycle of each goroutine and your entire program as a consequence.

Having seen bad program designs, from worst to less bad:

- Using `os.Exit(1)` to terminate the program: goroutines do not terminate gracefully
- Hanging shutdowns when waiting for goroutines to complete
- Exiting all goroutines at the same time when cancelling a shared `context.Context`, when a shutdown order should be needed
- Waiting for all goroutines to finish using a single waitgroup `wg` with `wg.Wait()`
- Waiting on multiple `done` signal channels where one could block others from being canceled

I decided to write this library to ease the task in all `main.go`'s `main()` functions.

## Setup

```sh
go get github.com/qdm12/goshutdown
```

## Usage

### Example

This is a very simple example showing how to run two goroutines `badDeadlock` and `goodCleanup` where `badDeadlock` hangs when exiting and `goodCleanup` does some cleanup in 500ms.

We configure them to be shutdown in order, where `badDeadlock` should be shutdown first and `goodCleanup` after.

Our shutdown order is given a 3 seconds timeout, and each of our goroutine shutdown handlers use the default 1 second.

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/order"
)

func main() {
	const orderTimeout = 3 * time.Second
	orderSettings := order.Settings{
		Timeout:   orderTimeout,
		OnSuccess: func(name string) { log.Println(name + " terminated 🙌") },
		OnFailure: func(name string, err error) { log.Println(name + " did not terminate 😱: " + err.Error()) },
	}
	order := order.New("order", orderSettings)

	handlerA, ctxA, doneA := goroutine.New("badDeadlock", goroutine.Settings{})
	go badDeadlock(ctxA, doneA)
	order.Append(handlerA)

	handlerB, ctxB, doneB := goroutine.New("goodCleanup", goroutine.Settings{})
	go goodCleanup(ctxB, doneB)
	order.Append(handlerB)

	// do stuff, wait for OS signals etc.

	err := order.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func badDeadlock(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
	log.Println("😤 not exiting")
	theDeadLock := make(chan struct{})
	<-theDeadLock
}

func goodCleanup(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
	const ioTime = 500 * time.Millisecond
	log.Println("📤 doing some IO cleanup for " + ioTime.String())
	time.Sleep(ioTime)
}

```

The following is logged:

```
2021/06/09 15:13:55 😤 not exiting
2021/06/09 15:13:56 badDeadlock did not terminate 😱: goroutine shutdown timed out: after 1s
2021/06/09 15:13:56 📤 doing some IO cleanup for 500ms
2021/06/09 15:13:57 goodCleanup terminated 🙌
2021/06/09 15:13:57 ordered shutdown timed out: badDeadlock: goroutine shutdown timed out: after 1s
```

So what happened here?

1. The goroutine `badDeadlock` is shutdown using its context `ctxA`, but it hangs and never closes `doneA`
1. The shutdown logic waits 1 second (default timeout) for the `badDeadlock` goroutine to close `doneA`
1. It times out so it moves on to the next element to shutdown. Note you can set the `Critical: true` setting to `badDeadlock` to stop the order if it fails.
1. The goroutine `goodCleanup` is shutdown using its context `ctxB` and closes `doneB` after 500ms of fake cleanup.
1. Since it's within its 1 second timeout, it is terminated successfully
1. The order is now complete, returning an error since one of the elements timed out.

See the [examples](examples) for more examples.

### Available structures

- `goroutine.Handler` created using `goroutine.New("name", goroutine.Settings{})` for handling goroutines. This is the smallest piece in this `goshutdown`.
- `group.Handler` created using `group.New("name", group.Settings{})` for handling a group of handlers which will be shutdown **in parallel**.
- `order.Handler` created using `order.New("name", order.Settings{})` for handling an order of handlers which will be shutdown **sequentially**.

Each of these 3 handlers implement the [`handler.Handler`](handler/handler.go) interface:

```go
// Handler is the minimal common interface for shutdown items.
type Handler interface {
	// Name returns the name assigned to the handler.
	Name() string
	// IsCritical returns true if the shutdown process is critical and further
	// operations should be dropped it it cannot be done.
	IsCritical() bool
	// Shutdown initiates the shutdown process and returns an error if it fails.
	Shutdown(ctx context.Context) (err error)
}
```

Therefore they can also be nested within each other. For example you could have an order of 1 group handler, 1 goroutine handler and another group handler.

### Settings

Each handler (goroutine, group and order) has their own settings structure.

What is common:

- `Timeout`: the maximum time allowed to shutdown the handler
- `Critical`: is the handler critical when viewed by a parent handler? If it is set to true, a parent handler would stop the shutdown operations if it cannot be terminated.

What is available to `group.Handler` and `order.Handler` only:

- `onSuccess` is a function executing as soon as a child handler is successfully terminated. This can be useful for logging purposes for example.
- `onFailure` is a function executing as soon as a child handler is not terminated on time. This can be useful for logging purposes for example.

### Save on imports

If you feel like you have too many import statements for this library, you can just import `"github.com/qdm12/goshutdown"` which has functions and type aliases to the `goroutine`, `order` and `group` subpackages. 

For example:

```go
package main

import (
	"context"
	"log"

	"github.com/qdm12/goshutdown"
)

func main() {
	order := goshutdown.NewOrderHandler("order", goshutdown.OrderSettings{})

	handlerA, ctxA, doneA := goshutdown.NewGoRoutineHandler("functionA", goshutdown.GoRoutineSettings{})
	go functionA(ctxA, doneA)
	order.Append(handlerA)

	err := order.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func functionA(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
}

```

## Safety to use

- The code is fully test covered by unit and integrationt tests
- The code is linted using `golangci-lint` with almost all its linters activated
- It is already in use in multiple Go projects with thousands of users:
  - [gluetun](https://github.com/qdm12/gluetun)
- You can use generated mocks (with [github.com/golang/mock](https://github.com/golang/mock)) for your own tests with for example

    ```go
    import (
      "github.com/qdm12/goshutdown/order/mock_order"
    )
    ```

    Or use the shorter import path `"github.com/qdm12/goshutdown/mock"` which contains shorthand constructors for the mocks.

## Bug and feature request

- [Create an issue](https://github.com/qdm12/goshutdown/issues/new) or [a discussion](https://github.com/qdm12/goshutdown/discussions) for feature requests or bugs.

## Questions

- Rename `Group` to `Parallel`/`Wave`?
- Rename `Order` to `Sequential`?
- Is that shortapi pattern OK for [`shortapi.go`](shortapi.go) and [`mock/shortapi.go`](mock/shortapi.go)?
