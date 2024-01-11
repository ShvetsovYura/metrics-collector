package handlers

import (
	"fmt"

	"github.com/go-chi/chi"
)

func ServerRouter(m Storage) chi.Router {
	r := chi.NewRouter()
	r.Get("/", MetricGetCurrentValuesHandler(m))
	r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", metricType, metricName, metricValue), MetricUpdateHandler(m))
	r.Get(fmt.Sprintf("/value/{%s}/{%s}", metricType, metricName), MetricGetValueHandler(m))
	return r
}
