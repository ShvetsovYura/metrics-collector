package handlers

import (
	"github.com/ShvetsovYura/metrics-collector/internal/types"
	"github.com/go-chi/chi"
)

func ServerRouter(m types.Stored) chi.Router {
	r := chi.NewRouter()
	r.Get("/", MetricGetCurrentValuesHandler(m))
	r.Post("/update/{mType}/{mName}/{mVal}", MetricUpdateHandler(m))
	r.Get("/value/{mType}/{mName}", MetricGetValueHandler(m))
	return r
}
