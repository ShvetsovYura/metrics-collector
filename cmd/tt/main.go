package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	wg := sync.WaitGroup{}
	fmt.Println("start app")
	go upd(ctx, &wg)
	go send(ctx, &wg)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigChan

	wg.Wait()
}

func upd(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		select {
		case <-time.After(time.Second * 2):
			fmt.Println("get metrics")
		case <-ctx.Done():
			fmt.Println("upd def")
			wg.Done()
		}
	}
}

func send(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	for {

		select {
		case <-time.After(time.Second * 10):
			fmt.Println("send metrics")
		case <-ctx.Done():
			fmt.Println("send def")
			wg.Done()
		}
	}
}
