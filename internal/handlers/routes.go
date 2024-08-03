package handlers

import (
	"context"
	"fmt"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/middlewares"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
)

// StorageReader, интерфейс, определяющий поддержку чтение данных из стораджа.
type StorageReader interface {
	GetGauge(ctx context.Context, name string) (models.Gauge, error)
	GetCounter(ctx context.Context, name string) (models.Counter, error)
	Ping(ctx context.Context) error
	ToList(ctx context.Context) ([]string, error)
}

// StorageWriter, интерфейс, определяющий поддержку запись данных из сторадж.
type StorageWriter interface {
	SetGauge(ctx context.Context, name string, val float64) error
	SetCounter(ctx context.Context, name string, val int64) error
	SaveGaugesBatch(context.Context, map[string]models.Gauge) error
	SaveCountersBatch(context.Context, map[string]models.Counter) error
}

// Storage, интерфейс работы со стораджем.
type Storage interface {
	StorageReader
	StorageWriter
}

// ServerRouter, функция объявления роутинга http-запросов и их обработчиков.
func ServerRouter(s Storage, key string, privateKeyPath string) chi.Router {
	logger.NewHTTPLogger()

	r := chi.NewRouter()
	// if privateKeyPath != "" {
	// 	r.Use(middlewares.DecryptMessage(privateKeyPath))
	// }
	r.Use(middlewares.CheckRequestHashHeader(key))

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(httplog.RequestLogger(logger.HTTPLogger))
	r.Use(middlewares.WithUnzipRequest)
	r.Use(middlewares.ResposeHeaderWithHash(key))

	r.Get("/", MetricGetCurrentValuesHandler(s))

	pattern := fmt.Sprintf("/update/{%s}/{%s}/{%s}", internal.MetricTypePathParam, internal.MetricNamePathParam, internal.MetricValuePathParam)
	r.Post(pattern, MetricUpdateHandler(s))

	pattern = fmt.Sprintf("/value/{%s}/{%s}", internal.MetricTypePathParam, internal.MetricNamePathParam)
	r.Get(pattern, MetricGetValueHandler(s))

	r.Post("/update/", MetricUpdateHandlerWithBody(s))
	r.Post("/updates/", MetricBatchUpdateHandler(s))
	r.Post("/value/", MetricGetValueHandlerWithBody(s))
	r.Get("/ping", DBPingHandler(s))

	r.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/cmdline", pprof.Handler("cmdline").ServeHTTP)
		r.Get("/profile", pprof.Handler("profile").ServeHTTP)
		r.Get("/symbol", pprof.Handler("symbol").ServeHTTP)
		r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.Get("/heap", pprof.Handler("heap").ServeHTTP)
		r.Get("/trace", pprof.Trace)
	})

	return r
}
