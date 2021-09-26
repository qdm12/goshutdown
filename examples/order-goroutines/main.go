package main

import (
	"context"
	"log"
	"time"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/order"
)

func main() {
	order := order.New("order",
		order.OptionTimeout(time.Second),
		order.OptionOnSuccess(func(name string) { log.Println(name + " terminated 🙌") }),
		order.OptionOnFailure(func(name string, err error) { log.Println(name + " did not terminate 😱: " + err.Error()) }),
	)

	handlerA, ctxA, doneA := goroutine.New("functionA")
	go functionA(ctxA, doneA)
	order.Append(handlerA)

	handlerB, ctxB, doneB := goroutine.New("functionB")
	go functionB(ctxB, doneB)
	order.Append(handlerB)

	handlerC, ctxC, doneC := goroutine.New("functionC")
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
