package handlers

import (
	"context"
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/middlewares"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

type Storage interface {
	SetGauge(ctx context.Context, name string, val float64) error
	SetCounter(ctx context.Context, name string, val int64) error
	GetGauge(ctx context.Context, name string) (metric.Gauge, error)
	GetCounter(ctx context.Context, name string) (metric.Counter, error)
	Ping(ctx context.Context) error
	ToList(ctx context.Context) ([]string, error)
	Save() error
	Restore(context.Context) error
	SaveGaugesBatch(context.Context, map[string]metric.Gauge) error
	SaveCountersBatch(context.Context, map[string]metric.Counter) error
}

func ServerRouter(s Storage, key string) chi.Router {

	logger.NewHTTPLogger()

	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(httplog.RequestLogger(logger.HTTPLogger))
	if key != "" {
		r.Use(middlewares.HashCheck)
	}

	r.Get("/", middlewares.WithUnzipRequest(MetricGetCurrentValuesHandler(s)))
	r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", internal.MetricTypePathParam, internal.MetricNamePathParam, internal.MetricValuePathParam), middlewares.WithUnzipRequest(MetricUpdateHandler(s)))
	r.Get(fmt.Sprintf("/value/{%s}/{%s}", internal.MetricTypePathParam, internal.MetricNamePathParam), middlewares.WithUnzipRequest(MetricGetValueHandler(s)))
	r.Post("/update/", middlewares.WithUnzipRequest(MetricUpdateHandlerWithBody(s)))
	r.Post("/updates/", middlewares.WithUnzipRequest(MetricBatchUpdateHandler(s)))
	r.Post("/value/", middlewares.WithUnzipRequest(MetricGetValueHandlerWithBody(s)))
	r.Get("/ping", DBPingHandler(s))

	return r
}
