package agent

import (
	"encoding/json"
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
)

type Sender interface {
	Send(string, string)
}
type metrics map[string]Sender
type gauge float64
type counter int64

func NewMetrics(initSize int) metrics {
	m := make(map[string]Sender, initSize)
	m[counterMetricFieldName] = counter(0)
	return m
}

func (g gauge) Send(mName string, baseURL string) {
	link := fmt.Sprintf("http://%s/update/", baseURL)
	val := float64(g)
	obj := models.Metrics{
		ID:    mName,
		MType: "gauge",
		Value: &val,
	}
	data, _ := json.Marshal(obj)

	sendMetric(data, link, "application/json")
}

func (c counter) Send(mName string, baseURL string) {
	link := fmt.Sprintf("http://%s/update/", baseURL)
	val := int64(c)
	data, _ := json.Marshal(models.Metrics{
		ID:    counterMetricFieldName,
		MType: internal.InCounterName,
		Delta: &val,
	})
	sendMetric(data, link, "application/json")

}
