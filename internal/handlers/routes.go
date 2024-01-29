package handlers

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func ServerRouter(s Storage) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "application/json", "text/html"))

	r.Get("/", middlewares.WithLog(MetricGetCurrentValuesHandler(s)))
	r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", metricType, metricName, metricValue), middlewares.WithLog(MetricUpdateHandler(s)))
	r.Get(fmt.Sprintf("/value/{%s}/{%s}", metricType, metricName), middlewares.WithLog(MetricGetValueHandler(s)))
	r.Post("/update/", middlewares.WithLog(MetricUpdateHandlerWithBody(s)))
	r.Post("/value/", middlewares.WithLog(MetricGetValueHandlerWithBody(s)))
	return r
}
