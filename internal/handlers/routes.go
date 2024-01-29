package handlers

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func ServerRouter(s Storage) chi.Router {

	logger.NewHTTPLogger()

	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "application/json", "text/html"))

	r.Use(httplog.RequestLogger(logger.HTTPLogger))

	r.Get("/", MetricGetCurrentValuesHandler(s))
	r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", metricType, metricName, metricValue), MetricUpdateHandler(s))
	r.Get(fmt.Sprintf("/value/{%s}/{%s}", metricType, metricName), MetricGetValueHandler(s))
	r.Post("/update/", MetricUpdateHandlerWithBody(s))
	r.Post("/value/", MetricGetValueHandlerWithBody(s))
	return r
}
