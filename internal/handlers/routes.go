package handlers

import (
	"fmt"

	"github.com/go-chi/chi"
)

func ServerRouter(s Storage) chi.Router {
	r := chi.NewRouter()
	r.Get("/", MetricGetCurrentValuesHandler(s))
	r.Post(fmt.Sprintf("/update/{%s}/{%s}/{%s}", metricType, metricName, metricValue), MetricUpdateHandler(s))
	r.Get(fmt.Sprintf("/value/{%s}/{%s}", metricType, metricName), MetricGetValueHandler(s))
	return r
}
