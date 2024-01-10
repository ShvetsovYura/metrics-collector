package handlers

import (
	"github.com/go-chi/chi"
)

func ServerRouter(m Storage) chi.Router {
	r := chi.NewRouter()
	r.Get("/", MetricGetCurrentValuesHandler(m))
	r.Post("/update/{mType}/{mName}/{mVal}", MetricUpdateHandler(m))
	r.Get("/value/{mType}/{mName}", MetricGetValueHandler(m))
	return r
}
