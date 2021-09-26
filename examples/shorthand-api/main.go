package main

import (
	"context"
	"log"

	"github.com/qdm12/goshutdown"
)

func main() {
	order := goshutdown.NewOrderHandler("order")

	handlerA, ctxA, doneA := goshutdown.NewGoRoutineHandler("functionA")
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
