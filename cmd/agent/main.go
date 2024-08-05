// Запускает агента по сбору метрик.

package main

import (
	"context"
	"fmt"
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

	opts := agent.ReadOptions()

	a := agent.NewAgent(metricsCount, opts)

	logger.Log.Info("Start AGENT app")
	showBuildInfo("Build version: ", buildVersion)
	showBuildInfo("Build date: ", buildDate)
	showBuildInfo("Build commit: ", buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)

	defer stop()
	a.Run(ctx)
}

func showBuildInfo(caption string, v string) {
	if v == "" {
		logger.Log.Infof("%s: N/A", caption)
	} else {
		logger.Log.Infof("%s: %s", caption, v)

	}
}
