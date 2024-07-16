// Запускает агента по сбору метрик.

package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

const metricsCount int = 40

func main() {
	err := logger.InitLogger("info")
	if err != nil {
		fmt.Println("Не удалось инициализировать лог")
	}

	opts := new(agent.Options)
	opts.ParseArgs()

	if err := opts.ParseEnvs(); err != nil {
		log.Fatal(err.Error())
	}

	a := agent.NewAgent(metricsCount, opts)

	logger.Log.Info("Start agent app")
	logger.Log.Infof("Build version: %s", buildVersion)
	logger.Log.Infof("Build date: %s", buildDate)
	logger.Log.Infof("Build commit: %s", buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)

	defer stop()
	a.Run(ctx)
}
