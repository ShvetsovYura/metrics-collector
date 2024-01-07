package main

import (
	"fmt"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"github.com/go-chi/chi"
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
	r := chi.NewRouter()
	r.Get("/", handlers.MetricGetCurrentValuesHandler(&m))
	r.Post("/update/{mType}/{mName}/{mVal}", handlers.MetricUpdateHandler(&m))
	r.Get("/value/{mType}/{mName}", handlers.MetricGetValueHandler(&m))
	return http.ListenAndServe(opts.EndpointAddr, r)
}
