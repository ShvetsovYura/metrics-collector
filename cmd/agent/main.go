// Запускает агента по сбору метрик.

package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	httpclient "github.com/ShvetsovYura/metrics-collector/internal/agent/http_client"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

const metricsCount int = 40

func main() {
	fmt.Println("Запускается АГЕНТ сбора метрик...")
	opts := agent.ReadOptions()
	err := logger.InitLogger(opts.LogLevel)
	if err != nil {
		fmt.Println("Не удалось инициализировать лог")
	}

	metricSender := httpclient.NewMetricSender(
		"http://"+opts.EndpointAddr+"/update/", agent.DefaultContentType, opts.Key, opts.CryptoKey,
	)
	metricCollection := agent.NewMetricCollector(metricsCount)
	a := agent.NewAgent(metricCollection, metricSender, opts)
	showBuildInfo("Build version: ", buildVersion)
	showBuildInfo("Build date: ", buildDate)
	showBuildInfo("Build commit: ", buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()
	a.Run(ctx)
	logger.Log.Info("работа АГЕНТА сбора метрик завершена")
}

func showBuildInfo(caption string, v string) {
	if v == "" {
		logger.Log.Infof("%s: N/A", caption)
	} else {
		logger.Log.Infof("%s: %s", caption, v)

	}
}
