// Запускает web-сервер по сбору/обработке метрик.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func serverFactory(serverType string, hashKey string, trustedSubnet string) (server.IServer, error) {
	if serverType == "http" {
		return server.NewHTTPServer(), nil
	}
	if serverType == "grpc" {
		return server.NewGRPCServer(trustedSubnet, hashKey), nil
	}
	return nil, errors.New("не удалось определить тип запускаемого сервера")
}

func main() {
	opts := server.ReadOptions()
	err := logger.InitLogger(opts.LogLevel)
	if err != nil {
		fmt.Printf("Не удалось инициализировать лог, %s \n", err.Error())
	}
	logger.Log.Info(*opts)

	serverType, err := serverFactory(opts.ServerType, opts.Key, opts.TrustedSubnet)
	if err != nil {
		log.Fatal(err.Error())
	}
	srv := server.NewServer(serverType, 40, opts)

	logger.Log.Infof("Start server with options: %v", *opts)
	showBuildInfo("Build version: ", buildVersion)
	showBuildInfo("Build date: ", buildDate)
	showBuildInfo("Build commit: ", buildCommit)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

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
