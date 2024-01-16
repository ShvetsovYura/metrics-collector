package handlers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"github.com/go-chi/chi"
)

type Storage interface {
	UpdateGauge(name string, val float64) error
	UpdateCounter(val int64)
	GetGauge(name string) (metric.Gauge, error)
	GetCounter() (metric.Counter, error)
	ToList() []string
}

type Store interface {
	SetGauge(name string, val float64) error
	SetCounter() error
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
				m.UpdateCounter(parsedVal)
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

		if mName == gaugeName {
			v, err := m.GetGauge(mName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, v.ToString())
		} else if mName == counterName {
			v, err := m.GetCounter()
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, v.ToString())
		}

	}
}

func MetricGetCurrentValuesHandler(m Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, strings.Join(m.ToList(), ", "))
	}
}
