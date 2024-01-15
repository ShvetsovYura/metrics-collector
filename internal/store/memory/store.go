package memory

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/store/metric"
)

type MemStore struct {
	metrics map[string]metric.Metric
}

func NewStore(metricsCount int) *MemStore {
	ms := new(MemStore)
	ms.metrics = make(map[string]metric.Metric, metricsCount)
	ms.metrics["PollCounter"] = metric.Counter(0)
	return ms
}

func (m *MemStore) SetGauge(name string, val float64) error {
	m.metrics[name] = metric.Gauge(val)
	return nil
}

func (m *MemStore) SetCounter() error {
	val := m.metrics["PollCounter"]
	m.metrics["PollCounter"] = val.(metric.Counter) + metric.Counter(1)
	return nil
}

func (m *MemStore) GetVal(name string) (metric.Metric, error) {
	if val, ok := m.metrics[name]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStore) ToList() []string {
	var list []string

	for _, c := range m.metrics {
		list = append(list, c.ToString())
	}
	return list
}
