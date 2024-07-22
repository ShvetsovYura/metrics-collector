// Запускает web-сервер по сбору/обработке метрик.

package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	err := logger.InitLogger("info")
	if err != nil {
		fmt.Println("Не удалось инициализировать лог")
	}

	opts := new(server.Options)
	opts.ParseArgs()

	if err := opts.ParseEnvs(); err != nil {
		logger.Log.Fatal(err.Error())
	}

	logger.Log.Info(*opts)

	srv := server.NewServer(40, opts)

	logger.Log.Infof("Start server with options: %v", *opts)
	showBuildInfo("Build version: ", buildVersion)
	showBuildInfo("Build date: ", buildDate)
	showBuildInfo("Build commit: ", buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)

	defer stop()

	if err := srv.Run(ctx); err != nil {
		panic(err)
	}
}

func showBuildInfo(caption string, v string) {
	if v == "" {
		logger.Log.Infof("%s: N/A", caption)
	} else {
		logger.Log.Infof("%s: %s", caption, v)

	}
}
