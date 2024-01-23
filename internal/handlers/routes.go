package handlers

import (
	"github.com/ShvetsovYura/metrics-collector/internal/middlewares"
	"github.com/go-chi/chi"
)

func ServerRouter(s Storage) chi.Router {
	r := chi.NewRouter()
	r.Get("/", middlewares.WithLog(MetricGetCurrentValuesHandler(s)))
	r.Post("/update/", middlewares.WithLog(MetricUpdateHandlerWithBody(s)))
	r.Post("/value/", middlewares.WithLog(MetricGetValueHandlerWithBody(s)))
	// r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", metricType, metricName, metricValue), middlewares.WithLog(MetricUpdateHandler(s)))
	// r.Get(fmt.Sprintf("/value/{%s}/{%s}", metricType, metricName), middlewares.WithLog(MetricGetValueHandler(s)))
	return r
}
