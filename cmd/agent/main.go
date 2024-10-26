// Запускает агента по сбору метрик.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	grpcclient "github.com/ShvetsovYura/metrics-collector/internal/agent/grpc_client"
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

	metricCollection := agent.NewMetricCollector(metricsCount)
	client, err := selectSenderClient(opts.ClientType, opts.EndpointAddr, "", opts.Key, opts.CryptoKey)
	if err != nil {
		log.Fatal(err)
	}
	a := agent.NewAgent(metricCollection, client, opts)
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

func selectSenderClient(clientType string, addr string, contentType string, hashKey string, cryptoKey string) (agent.Sender, error) {
	switch clientType {
	case "http":
		return httpclient.NewClient("http://"+addr+"/update/", contentType, hashKey, cryptoKey), nil
	case "grpc":
		return grpcclient.NewClient(addr, hashKey)
	default:
		return nil, errors.New("не найден указанный тип клиента")
	}

}
