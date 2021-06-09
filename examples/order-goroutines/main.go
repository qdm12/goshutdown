package main

import (
	"context"
	"log"
	"time"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/order"
)

func main() {
	const orderTimeout = time.Second
	orderSettings := order.Settings{
		Timeout: orderTimeout,
		OnSuccess: func(name string) {
			log.Println(name + " terminated 🙌")
		},
		OnFailure: func(name string, err error) {
			log.Println(name + " did not terminate 😱: " + err.Error())
		},
	}
	order := order.New("order", orderSettings)

	goroutineSettings := goroutine.Settings{}

	handlerA, ctxA, doneA := goroutine.New("functionA", goroutineSettings)
	go functionA(ctxA, doneA)
	order.Append(handlerA)

	handlerB, ctxB, doneB := goroutine.New("functionB", goroutineSettings)
	go functionB(ctxB, doneB)
	order.Append(handlerB)

	handlerC, ctxC, doneC := goroutine.New("functionC", goroutineSettings)
	go functionC(ctxC, doneC)
	order.Append(handlerC)

	err := order.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func functionA(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
	log.Println("🔌 exiting on time!")
}

func functionB(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
	const ioTime = 500 * time.Millisecond
	log.Println("📤 doing some IO cleanup for " + ioTime.String())
	time.Sleep(ioTime)
}

func functionC(ctx context.Context, done chan<- struct{}) {
	defer close(done)
	<-ctx.Done()
	log.Println("😤 not exiting")
	theDeadLock := make(chan struct{})
	<-theDeadLock
}
