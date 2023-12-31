package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

func contains(s []string, val string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}

func run() error {
	r := chi.NewRouter()
	r.Post("/update/{mType}/{mName}/{mVal}", metricUpdateHandle)
	r.Get("/value/{mType}/{mName}", metricValueHandle)
	return http.ListenAndServe(`:8080`, r)
}

func metricUpdateHandle(w http.ResponseWriter, r *http.Request) {
	mType := chi.URLParam(r, "mType")
	mName := chi.URLParam(r, "mName")
	mVal := chi.URLParam(r, "mVal")

	if !contains(allowMetricTypes, mType) {
		w.WriteHeader(http.StatusBadRequest)
	}
	if mType == "gauge" {
		_, err := strconv.ParseFloat(mVal, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	}
	if mType == "counter" {
		_, err := strconv.ParseInt(mVal, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	w.WriteHeader(http.StatusOK)
}
