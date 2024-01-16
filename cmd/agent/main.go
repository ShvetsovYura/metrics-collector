package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
)

func main() {

	opts := new(agent.AgentOptions)
	opts.ParseArgs()
	if err := opts.ParseEnvs(); err != nil {
		log.Fatal(err.Error())
	}
	storage := memory.NewStorage(40)
	a := agent.NewAgent(storage, opts)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	wg := sync.WaitGroup{}
	fmt.Println("start app")
	a.Run(ctx, &wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-sigChan

	wg.Wait()
}
