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
