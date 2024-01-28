package handlers

import (
	"github.com/ShvetsovYura/metrics-collector/internal/middlewares"
	"github.com/go-chi/chi"
)

func ServerRouter(s Storage) chi.Router {
	r := chi.NewRouter()
	r.Get("/", middlewares.WithLog(middlewares.WithGzip(MetricGetCurrentValuesHandler(s))))
	r.Post("/update/", middlewares.WithLog(middlewares.WithGzip(MetricUpdateHandlerWithBody(s))))
	r.Post("/value/", middlewares.WithLog(middlewares.WithGzip(MetricGetValueHandlerWithBody(s))))
	return r
}
