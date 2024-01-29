package handlers

import (
	"bytes"
	"encoding/json"
	"strconv"

	"io"
	"net/http"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"github.com/go-chi/chi/v5"
)

type Storage interface {
	SetGauge(name string, val float64) error
	SetCounter(name string, val int64) error
	GetGauge(name string) (metric.Gauge, error)
	GetCounter(name string) (metric.Counter, error)
	ToList() []string
}

type Store interface {
	SetGauge(name string, val float64) error
	SetCounter() error
}

func MetricUpdateHandler(m Storage) http.HandlerFunc {
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
				m.SetGauge(mName, parsedVal)
			}

		}
		if mType == counterName {
			parsedVal, err := strconv.ParseInt(mVal, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				m.SetCounter(mName, parsedVal)
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func MetricGetValueHandler(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, metricType)
		mName := chi.URLParam(r, metricName)

		if !util.Contains([]string{gaugeName, counterName}, mType) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if mType == gaugeName {
			v, err := m.GetGauge(mName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, v.ToString())
		} else if mType == counterName {

			v, err := m.GetCounter(mName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, v.ToString())
		}

	}

}

func MetricUpdateHandlerWithBody(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity := models.Metrics{}

		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(b, &entity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !util.Contains([]string{internal.InGaugeName, internal.InCounterName}, entity.MType) {
			w.WriteHeader(http.StatusBadRequest)
		}

		var marshalVal []byte
		var marshalErr error
		if entity.MType == internal.InGaugeName {
			m.SetGauge(entity.ID, *entity.Value)
			val, _ := m.GetGauge(entity.ID)
			actualVal := models.Metrics{
				ID:    entity.ID,
				MType: internal.InGaugeName,
				Value: val.GetRawValue(),
			}
			marshalVal, marshalErr = json.Marshal(actualVal)

		} else if entity.MType == internal.InCounterName {
			m.SetCounter(entity.ID, *entity.Delta)
			val, _ := m.GetCounter(entity.ID)
			actualVal := models.Metrics{
				ID:    entity.ID,
				MType: internal.InCounterName,
				Delta: val.GetRawValue(),
			}
			marshalVal, marshalErr = json.Marshal(actualVal)
		}

		if marshalErr != nil {
			http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(marshalVal)
	}
}

func MetricGetValueHandlerWithBody(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		entity := models.Metrics{}
		var answer []byte

		_, err := buf.ReadFrom(r.Body)
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if err := json.Unmarshal(buf.Bytes(), &entity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if !util.Contains([]string{internal.InGaugeName, internal.InCounterName}, entity.MType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if entity.MType == internal.InGaugeName {
			v, err := m.GetGauge(entity.ID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			val, err := json.Marshal(models.Metrics{
				ID:    entity.ID,
				MType: internal.InGaugeName,
				Value: v.GetRawValue(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			answer = val
		} else if entity.MType == internal.InCounterName {
			v, _ := m.GetCounter(entity.ID)
			val, err := json.Marshal(models.Metrics{
				ID:    entity.ID,
				MType: internal.InCounterName,
				Delta: v.GetRawValue(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			answer = val
		}

		w.WriteHeader(http.StatusOK)
		w.Write(answer)
	}

}

func MetricGetCurrentValuesHandler(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, strings.Join(m.ToList(), ", "))
	}
}
