package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io"
	"net/http"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/ShvetsovYura/metrics-collector/internal/types"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
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

// func MetricUpdateHandler(m Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		mType := chi.URLParam(r, metricType)
// 		mName := chi.URLParam(r, metricName)
// 		mVal := chi.URLParam(r, metricValue)

// 		if !util.Contains([]string{gaugeName, counterName}, mType) {
// 			w.WriteHeader(http.StatusBadRequest)
// 		}
// 		if mType == gaugeName {
// 			parsedVal, err := strconv.ParseFloat(mVal, 64)
// 			if err != nil {
// 				w.WriteHeader(http.StatusBadRequest)
// 			} else {
// 				m.UpdateGauge(mName, parsedVal)
// 			}

// 		}
// 		if mType == counterName {
// 			parsedVal, err := strconv.ParseInt(mVal, 10, 64)
// 			if err != nil {
// 				w.WriteHeader(http.StatusBadRequest)
// 			} else {
// 				m.UpdateCounter(parsedVal)
// 			}
// 		}

// 		w.WriteHeader(http.StatusOK)
// 	}
// }

func MetricUpdateHandlerWithBody(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var buf bytes.Buffer
		entity := types.Metrics{}

		// _, err := buf.ReadFrom(r.Body)
		b, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(b, &entity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if !util.Contains([]string{gaugeName, counterName}, entity.MType) {
			w.WriteHeader(http.StatusBadRequest)
		}
		if entity.MType == gaugeName {
			m.UpdateGauge(entity.ID, *entity.Value)

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// func MetricGetValueHandler(m Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		mType := chi.URLParam(r, metricType)
// 		mName := chi.URLParam(r, metricName)

// 		if !util.Contains([]string{gaugeName, counterName}, mType) {
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}

// 		if mName == gaugeName {
// 			v, err := m.GetGauge(mName)
// 			if err != nil {
// 				w.WriteHeader(http.StatusNotFound)
// 				return
// 			}
// 			w.WriteHeader(http.StatusOK)
// 			io.WriteString(w, v.ToString())
// 		} else if mName == counterName {
// 			v, err := m.GetCounter()
// 			if err != nil {
// 				w.WriteHeader(http.StatusNotFound)
// 				return
// 			}
// 			w.WriteHeader(http.StatusOK)
// 			io.WriteString(w, v.ToString())
// 		}

// 	}

// }

func MetricGetValueHandlerWithBody(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		entity := types.Metrics{}
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
			v, err := m.GetGauge(entity.ID)

			if err != nil {
				fmt.Println(entity)
				fmt.Println(err.Error())
				val, _ := json.Marshal(types.Metrics{
					ID:    entity.ID,
					MType: "gauge",
					Value: v.GetRawValue(),
				})
				w.WriteHeader(http.StatusOK)
				w.Write((val))
				return
			}

			val, err := json.Marshal(types.Metrics{
				ID:    entity.ID,
				MType: "gauge",
				Value: v.GetRawValue(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(val)
		} else if entity.MType == counterName {
			v, err := m.GetCounter()
			if err != nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}
			val, err := json.Marshal(types.Metrics{
				ID:    "PollCounter",
				MType: "counter",
				Delta: v.GetRawValue(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(val)
		}

	}

}

func MetricGetCurrentValuesHandler(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, strings.Join(m.ToList(), ", "))
	}
}
