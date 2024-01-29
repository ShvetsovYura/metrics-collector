package handlers

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func ServerRouter(s Storage) chi.Router {

	logger.NewHTTPLogger()

	r := chi.NewRouter()

	r.Use(middleware.Compress(5))
	r.Use(httplog.RequestLogger(logger.HTTPLogger))

	// Я так и не нашел, как можно с помощью мидлварь chi декодировать запрос от агента
	// поэтому оставил свою реализацию только для чтения gzip
	r.Get("/", middlewares.WithUnzipRequest(MetricGetCurrentValuesHandler(s)))
	r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", metricType, metricName, metricValue), middlewares.WithUnzipRequest(MetricUpdateHandler(s)))
	r.Get(fmt.Sprintf("/value/{%s}/{%s}", metricType, metricName), middlewares.WithUnzipRequest(MetricGetValueHandler(s)))
	r.Post("/update/", middlewares.WithUnzipRequest(MetricUpdateHandlerWithBody(s)))
	r.Post("/value/", middlewares.WithUnzipRequest(MetricGetValueHandlerWithBody(s)))
	return r
}
