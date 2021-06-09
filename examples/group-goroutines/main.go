package main

import (
	"context"
	"log"
	"time"

	"github.com/qdm12/goshutdown/goroutine"
	"github.com/qdm12/goshutdown/group"
)

func main() {
	const groupTimeout = time.Second
	groupSettings := group.Settings{
		Timeout: groupTimeout,
		OnSuccess: func(name string) {
			log.Println(name + " terminated 🙌")
		},
		OnFailure: func(name string, err error) {
			log.Println(name + " did not terminate 😱: " + err.Error())
		},
	}
	group := group.New("group", groupSettings)

	goroutineSettings := goroutine.Settings{}

	handlerA, ctxA, doneA := goroutine.New("functionA", goroutineSettings)
	go functionA(ctxA, doneA)
	group.Add(handlerA)

	handlerB, ctxB, doneB := goroutine.New("functionB", goroutineSettings)
	go functionB(ctxB, doneB)
	group.Add(handlerB)

	handlerC, ctxC, doneC := goroutine.New("functionC", goroutineSettings)
	go functionC(ctxC, doneC)
	group.Add(handlerC)

	err := group.Shutdown(context.Background())
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
