// Запускает web-сервер по сбору/обработке метрик.

package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/server"
)

func main() {
	logger.InitLogger("info")

	opts := new(server.ServerOptions)
	opts.ParseArgs()

	if err := opts.ParseEnvs(); err != nil {
		logger.Log.Fatal(err.Error())
	}
	logger.Log.Info(*opts)
	srv := server.NewServer(40, opts)
	logger.Log.Infof("Start server with options: %v", *opts)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	if err := srv.Run(ctx); err != nil {
		panic(err)
	}
}
