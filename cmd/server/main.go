package main

import (
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
)

type gauge float64
type counter int64

var allowMetricTypes = []string{`gauge`, `counter`}

type MemStorage struct {
	gauges   map[string]gauge
	counters map[string]counter
}

type Stored interface {
	UpdateGauge(name string, val gauge)
	AddCounter(name string, val counter)
	Get(name string) any
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
