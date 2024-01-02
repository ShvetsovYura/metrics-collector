package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func MetricUpdateHandle(w http.ResponseWriter, r *http.Request) {
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
