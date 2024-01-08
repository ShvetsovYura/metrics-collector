package main

import (
	"fmt"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

func main() {
	opts := new(util.ServerOptions)
	opts.ParseArgs()
	opts.ParseEnvs()

	fmt.Println(opts)
	if err := run(opts); err != nil {
		panic(err)
	}
}

func run(opts *util.ServerOptions) error {
	m := storage.New()
	router := handlers.ServerRouter(m)
	return http.ListenAndServe(opts.EndpointAddr, router)
}
