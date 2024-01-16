package main

import (
	"log"

	"github.com/ShvetsovYura/metrics-collector/internal/server"
)

func main() {
	opts := new(server.ServerOptions)
	opts.ParseArgs()

	if err := opts.ParseEnvs(); err != nil {
		log.Fatal(err.Error())
	}
	srv := server.NewServer(opts)
	log.Printf("Start server with options: %v", &opts)
	if err := srv.Run(); err != nil {
		panic(err)
	}
}
