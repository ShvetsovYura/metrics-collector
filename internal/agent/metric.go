package agent

import (
	"encoding/json"
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
)

type Sender interface {
	Send(string, string, string)
	MarshalToJSON(string) []byte
	GetObj(string) models.Metrics
}
type metrics map[string]Sender
type gauge float64
type counter int64

func NewMetrics(initSize int) metrics {
	m := make(map[string]Sender, initSize)
	m[counterMetricFieldName] = counter(0)
	return m
}

func (m metrics) SendBatch(baseURL string, key string) error {
	metricsBatch := make([]models.Metrics, 0, 100)
	for k, v := range m {
		metricsBatch = append(metricsBatch, v.GetObj(k))
	}
	logger.Log.Info(metricsBatch)
	link := fmt.Sprintf("http://%s/updates/", baseURL)
	data, err := json.Marshal(metricsBatch)
	if err != nil {
		return err
	}
	sendMetric(data, link, "application/json", key)
	return err
}

func (g gauge) Send(mName string, baseURL string, key string) {
	link := fmt.Sprintf("http://%s/update/", baseURL)

	data := g.MarshalToJSON(mName)
	sendMetric(data, link, "application/json", key)
}

func (g gauge) MarshalToJSON(mName string) []byte {
	data, _ := json.Marshal(g.GetObj(mName))
	return data
}

func (c counter) Send(mName string, baseURL string, key string) {
	link := fmt.Sprintf("http://%s/update/", baseURL)
	data := c.MarshalToJSON(mName)
	sendMetric(data, link, "application/json", key)
}

func (c counter) MarshalToJSON(mName string) []byte {
	data, _ := json.Marshal(c.GetObj(mName))
	return data
}

func (c counter) GetObj(mName string) models.Metrics {
	val := int64(c)

	return models.Metrics{
		ID:    mName,
		MType: internal.InCounterName,
		Delta: &val,
	}
}

func (g gauge) GetObj(mName string) models.Metrics {
	val := float64(g)

	return models.Metrics{
		ID:    mName,
		MType: internal.InGaugeName,
		Value: &val,
	}
}
