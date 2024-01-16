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

	go updateMetrics(ctx, &wg)
	go sendMetrics(ctx, &wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigChan

	wg.Wait()
}

func updateMetrics(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		select {
		case <-time.After(time.Second * 2):
			fmt.Println("update")
		case <-ctx.Done():
			fmt.Println("upd def")
			wg.Done()
		}
	}
}

func sendMetrics(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	for {

		select {
		case <-time.After(time.Second * 10):
			fmt.Println("send")
		case <-ctx.Done():
			fmt.Println("send def")
			wg.Done()
		}
	}
}
