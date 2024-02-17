package handlers

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func ServerRouter(s Storage) chi.Router {

	logger.NewHTTPLogger()

	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(httplog.RequestLogger(logger.HTTPLogger))

	r.Get("/", middlewares.WithUnzipRequest(MetricGetCurrentValuesHandler(s)))
	r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", internal.MetricTypePathParam, internal.MetricNamePathParam, internal.MetricValuePathParam), middlewares.WithUnzipRequest(MetricUpdateHandler(s)))
	r.Get(fmt.Sprintf("/value/{%s}/{%s}", internal.MetricTypePathParam, internal.MetricNamePathParam), middlewares.WithUnzipRequest(MetricGetValueHandler(s)))
	r.Post("/update/", middlewares.WithUnzipRequest(MetricUpdateHandlerWithBody(s)))
	r.Post("/updates/", middlewares.WithUnzipRequest(MetricBatchUpdateHandler(s)))
	r.Post("/value/", middlewares.WithUnzipRequest(MetricGetValueHandlerWithBody(s)))
	r.Get("/ping", DBPingHandler(s))

	return r
}
