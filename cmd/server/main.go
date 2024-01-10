package main

import (
	"log"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
)

const metiricsCount int = 40

func main() {
	opts := new(ServerOptions)
	opts.ParseArgs()

	if err := opts.ParseEnvs(); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Start server with options: %v", &opts)
	if err := run(opts); err != nil {
		panic(err)
	}
}

func run(opts *ServerOptions) error {
	m := storage.NewStorage(metiricsCount)
	router := handlers.ServerRouter(m)
	return http.ListenAndServe(opts.EndpointAddr, router)
}
