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

const (
	metricType  string = "mType"
	metricName  string = "mName"
	metricValue string = "mVal"
	gaugeName   string = "gauge"
	counterName string = "counter"
)

type Storage interface {
	UpdateGauge(name string, val float64)
	UpdateCounter(name string, val int64)
	GetVal(name string) (storage.Metric, error)
	ToList() []string
}

func MetricUpdateHandler(m Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, metricType)
		mName := chi.URLParam(r, metricName)
		mVal := chi.URLParam(r, metricValue)

		if !util.Contains([]string{gaugeName, counterName}, mType) {
			w.WriteHeader(http.StatusBadRequest)
		}
		if mType == gaugeName {
			parsedVal, err := strconv.ParseFloat(mVal, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				m.UpdateGauge(mName, parsedVal)
			}

		}
		if mType == counterName {
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

func MetricGetValueHandler(m Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, metricType)
		mName := chi.URLParam(r, metricName)

		if !util.Contains([]string{gaugeName, counterName}, mType) {
			w.WriteHeader(http.StatusNotFound)
			return
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

func MetricGetCurrentValuesHandler(m Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, strings.Join(m.ToList(), ", "))
	}
}
