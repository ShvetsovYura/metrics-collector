package main

import (
	"net/http"
	"strconv"
	"strings"
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

func contains(s []string, val string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(metricHandler))
}

func metricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if pathParts[1] != "update" {
		w.WriteHeader(http.StatusNotFound)
	}
	parts := pathParts[2:]
	if len(parts) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	mType := parts[0]
	allowMetricTypes := []string{`gauge`, `counter`}
	if !contains(allowMetricTypes, mType) {
		w.WriteHeader(http.StatusBadRequest)
	}
	if mType == "gauge" {
		_, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	if mType == "counter" {
		_, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	w.WriteHeader(http.StatusOK)
}
