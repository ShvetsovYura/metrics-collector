package handlers

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/ShvetsovYura/metrics-collector/internal/types"
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
		entity := types.Metrics{}

		b, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(b, &entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !util.Contains([]string{gaugeName, counterName}, entity.MType) {
			w.WriteHeader(http.StatusBadRequest)
		}

		if entity.MType == gaugeName {
			m.UpdateGauge(entity.ID, *entity.Value)
			m.SaveNow()

		} else if entity.MType == counterName {
			m.UpdateCounter(*entity.Delta)
		}

		val, err := m.GetGauge(entity.ID)
		actualVal := types.Metrics{
			ID:    entity.ID,
			MType: gaugeName,
			Value: val.GetRawValue(),
		}

		resp, err := json.Marshal(actualVal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

func MetricGetValueHandlerWithBody(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		entity := types.Metrics{}
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

		if !util.Contains([]string{gaugeName, counterName}, entity.MType) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if entity.MType == gaugeName {
			v, _ := m.GetGauge(entity.ID)

			val, err := json.Marshal(types.Metrics{
				ID:    entity.ID,
				MType: gaugeName,
				Value: v.GetRawValue(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			answer = val
		} else if entity.MType == counterName {
			v, _ := m.GetCounter()
			val, err := json.Marshal(types.Metrics{
				ID:    "PollCounter",
				MType: counterName,
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
