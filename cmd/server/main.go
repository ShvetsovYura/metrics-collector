package main

import (
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
	"github.com/go-chi/chi"
)

// type Stored interface {
// 	Update(name string, val gauge)
// 	AddCounter(name string, val counter)
// 	Get(name string) any
// }

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	m := storage.New()
	r := chi.NewRouter()
	r.Get("/", handlers.MetricGetCurrentValuesHandler(&m))
	r.Post("/update/{mType}/{mName}/{mVal}", handlers.MetricUpdateHandler(&m))
	r.Get("/value/{mType}/{mName}", handlers.MetricGetValueHandler(&m))
	return http.ListenAndServe(`:8080`, r)
}
