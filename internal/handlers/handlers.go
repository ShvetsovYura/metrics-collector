package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal/utilu"
)

func MetricHandler(w http.ResponseWriter, r *http.Request) {
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
	if !utilu.Contains(allowMetricTypes, mType) {
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
