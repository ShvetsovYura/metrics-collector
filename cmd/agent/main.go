package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

const metricsCount int = 40

func main() {
	logger.InitLogger("info")
	opts := new(agent.AgentOptions)
	opts.ParseArgs()
	if err := opts.ParseEnvs(); err != nil {
		log.Fatal(err.Error())
	}

	a := agent.NewAgent(metricsCount, opts)

	logger.Log.Info("Start agent app")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	a.Run(ctx)

}
