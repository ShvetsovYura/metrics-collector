package handlers

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

type Storage interface {
	UpdateGauge(name string, val float64) error
	UpdateCounter(val int64) error
	GetGauge(name string) (metric.Gauge, error)
	GetCounter() (metric.Counter, error)
	SaveNow()
	ToList() []string
}

type Store interface {
	SetGauge(name string, val float64) error
	SetCounter() error
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
			m.UpdateGauge(entity.ID, *entity.Value)
			m.SaveNow()
			val, _ := m.GetGauge(entity.ID)
			actualVal := models.Metrics{
				ID:    entity.ID,
				MType: internal.InGaugeName,
				Value: val.GetRawValue(),
			}
			marshalVal, marshalErr = json.Marshal(actualVal)

		} else if entity.MType == internal.InCounterName {
			m.UpdateCounter(*entity.Delta)
			m.SaveNow()
			val, _ := m.GetCounter()
			actualVal := models.Metrics{
				ID:    entity.ID,
				MType: internal.InCounterName,
				Delta: val.GetRawValue(),
			}
			marshalVal, marshalErr = json.Marshal(actualVal)
		}

		if marshalErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			v, _ := m.GetCounter()
			val, err := json.Marshal(models.Metrics{
				ID:    internal.CounterMetricFieldName,
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
