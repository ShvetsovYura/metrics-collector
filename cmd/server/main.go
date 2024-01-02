package main

import (
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
)

type MemStorage struct {
	metrics map[string]string
}

type storage interface {
	add(int64)
	change(float64)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(
		`:8080`,
		http.HandlerFunc(handlers.MetricHandler),
	)
}
