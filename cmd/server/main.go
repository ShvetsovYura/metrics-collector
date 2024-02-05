package main

import (
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
	srv := server.NewServer(40, opts)
	logger.Log.Infof("Start server with options: %v", *opts)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
