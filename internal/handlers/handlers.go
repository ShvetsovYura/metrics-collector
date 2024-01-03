package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal/storage"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"github.com/go-chi/chi"
)

func MetricUpdateHandler(m *storage.MemStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		mVal := chi.URLParam(r, "mVal")

		if !util.Contains([]string{"gauge", "counter"}, mType) {
			w.WriteHeader(http.StatusBadRequest)
		}
		if mType == "gauge" {
			parsedVal, err := strconv.ParseFloat(mVal, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				m.UpdateGauge(mName, parsedVal)
			}

		}
		if mType == "counter" {
			parsedVal, err := strconv.ParseInt(mVal, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				m.UpdateCounter(mName, parsedVal)
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func MetricGetValueHandler(m *storage.MemStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")

		if !util.Contains([]string{"gauge", "counter"}, mType) {
			w.WriteHeader(http.StatusNotFound)
		}

		v, err := m.GetVal(mName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, v.ToString())
	}
}

func MetricGetCurrentValuesHandler(m *storage.MemStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, strings.Join(m.ToList(), ", "))
	}
}
