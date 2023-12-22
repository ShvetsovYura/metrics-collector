package main

import (
	"net/http"
	"strings"
)

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
	parts := pathParts[2:]
	if len(parts) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	mType := parts[0]
	allowMetircTypes := []string{`gauge`, `counter`}
	if !contains(allowMetircTypes, mType) {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}
