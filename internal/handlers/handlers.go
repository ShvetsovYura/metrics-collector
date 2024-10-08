package handlers

import (
	"bytes"
	"encoding/json"
	"strconv"

	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

// MetricUpdateHandler, обновляет
func MetricUpdateHandler(m StorageWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, internal.MetricTypePathParam)
		mName := chi.URLParam(r, internal.MetricNamePathParam)
		mVal := chi.URLParam(r, internal.MetricValuePathParam)
		ctx := r.Context()

		switch mType {
		case internal.InGaugeName:
			parsedVal, err := strconv.ParseFloat(mVal, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			err = m.SetGauge(ctx, mName, parsedVal)
			if err != nil {
				logger.Log.Errorf("Ошибка установки значения для gauge: %s, значение: %f. %s", mName, parsedVal, err.Error())
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		case internal.InCounterName:
			parsedVal, err := strconv.ParseInt(mVal, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			err = m.SetCounter(ctx, mName, parsedVal)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

		default:
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func MetricGetValueHandler(m StorageReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mName := chi.URLParam(r, internal.MetricNamePathParam)
		mType := chi.URLParam(r, internal.MetricTypePathParam)

		switch mType {
		case internal.InGaugeName:
			v, err := m.GetGauge(ctx, mName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			_, err = io.WriteString(w, v.ToString())
			if err != nil {
				logger.Log.Errorf("Ошибка записи ответа, %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusOK)

		case internal.InCounterName:
			v, err := m.GetCounter(ctx, mName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			_, err = io.WriteString(w, v.ToString())
			if err != nil {
				logger.Log.Errorf("Ошибка записи ответа, %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusNotFound)

			return
		}
	}
}

func MetricUpdateHandlerWithBody(m Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		e := &models.MetricItem{}

		b, readerErr := io.ReadAll(r.Body)
		if readerErr != nil {
			http.Error(w, readerErr.Error(), http.StatusBadRequest)

			return
		}

		defer func() {
			bodyCloseErr := r.Body.Close()
			if bodyCloseErr != nil {
				logger.Log.Errorf("Ошибка закрытия тела ответа, %s", bodyCloseErr.Error())
			}
		}()
		logger.Log.Infof("incoming value: %v", string(b))
		w.Header().Set("Content-Type", "application/json")

		if err := json.Unmarshal(b, &e); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if !util.Contains([]string{internal.InGaugeName, internal.InCounterName}, e.MType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var (
			marshalVal []byte
			marshalErr error
		)

		switch e.MType {
		case internal.InGaugeName:
			err := m.SetGauge(ctx, e.ID, *e.Value)
			if err != nil {
				logger.Log.Errorf("Ошибка установки значения метрики gauge, %s", err.Error())
			}

			val, _ := m.GetGauge(ctx, e.ID)
			actualVal := models.MetricItem{
				ID:    e.ID,
				MType: internal.InGaugeName,
				Value: val.GetRawValue(),
			}

			marshalVal, marshalErr = json.Marshal(actualVal)
			if marshalErr != nil {
				logger.Log.Error(err.Error())
			}

		case internal.InCounterName:
			setErr := m.SetCounter(ctx, e.ID, *e.Delta)
			if setErr != nil {
				logger.Log.Errorf("Ошибка установки значнеия в метрики, %s", setErr.Error())
			}

			val, _ := m.GetCounter(ctx, e.ID)
			actualVal := models.MetricItem{
				ID:    e.ID,
				MType: internal.InCounterName,
				Delta: val.GetRawValue(),
			}
			logger.Log.Infof("metric actual value %v", actualVal)

			marshalVal, marshalErr = json.Marshal(actualVal)
			if marshalErr != nil {
				logger.Log.Error(marshalErr.Error())
			}
		}

		if marshalErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, writeErr := w.Write(marshalVal)
		if writeErr != nil {
			logger.Log.Errorf("Ошибка записи ответа, %s", writeErr.Error())
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func MetricGetValueHandlerWithBody(m StorageReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var (
			buf    bytes.Buffer
			answer []byte
		)

		entity := models.MetricItem{}

		_, err := buf.ReadFrom(r.Body)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				logger.Log.Errorf("Ошибка при закрытии тела ответа, %s ", closeErr.Error())
			}
		}()

		if unmarshalErr := json.Unmarshal(buf.Bytes(), &entity); unmarshalErr != nil {
			http.Error(w, unmarshalErr.Error(), http.StatusBadRequest)
		}

		if !util.Contains([]string{internal.InGaugeName, internal.InCounterName}, entity.MType) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if entity.MType == internal.InGaugeName {
			v, getGaugeErr := m.GetGauge(ctx, entity.ID)
			if getGaugeErr != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			val, marshalErr := json.Marshal(models.MetricItem{
				ID:    entity.ID,
				MType: internal.InGaugeName,
				Value: v.GetRawValue(),
			})
			if marshalErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			answer = val
		} else if entity.MType == internal.InCounterName {
			v, _ := m.GetCounter(ctx, entity.ID)
			val, marshalErr := json.Marshal(models.MetricItem{
				ID:    entity.ID,
				MType: internal.InCounterName,
				Delta: v.GetRawValue(),
			})

			if marshalErr != nil {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			answer = val
		}

		_, err = w.Write(answer)
		if err != nil {
			logger.Log.Errorf("Ошибка записи ответа, %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func MetricGetCurrentValuesHandler(m StorageReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "text/html")

		mList, err := m.ToList(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = io.WriteString(w, strings.Join(mList, ", "))
		if err != nil {
			logger.Log.Errorf("Ошибка записи ответа, %s", err.Error())
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DBPingHandler(m StorageReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := m.Ping(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func MetricBatchUpdateHandler(m StorageWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metricModels []models.MetricItem

		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		defer func() {
			closeErr := r.Body.Close()
			if closeErr != nil {
				logger.Log.Errorf("Ошибка при закрытии тела ответа, %s ", closeErr.Error())
			}
		}()

		errParse := json.Unmarshal(body, &metricModels)
		if errParse != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var (
			gauges   = make(map[string]models.Gauge, 100)
			counters = make(map[string]models.Counter, 100)
		)

		for _, mdl := range metricModels {
			switch mdl.MType {
			case internal.InGaugeName:
				gauges[mdl.ID] = models.Gauge(*mdl.Value)
			case internal.InCounterName:
				logger.Log.Infof("name: %s, val:%v", mdl.ID, *mdl.Delta)

				if v, ok := counters[mdl.ID]; ok {
					counters[mdl.ID] = v + models.Counter(*mdl.Delta)
				} else {
					counters[mdl.ID] = models.Counter(*mdl.Delta)
				}
			}
		}

		err = m.SaveCountersBatch(ctx, counters)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = m.SaveGaugesBatch(ctx, gauges)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}
}
