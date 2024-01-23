package memory

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
)

type Metric interface {
	ToString() string
}

type MemStorage struct {
	gaugeMetrics  map[string]metric.Gauge
	counterMetric metric.Counter
	// mux           sync.Mutex
}

func NewStorage(metricsCount int) *MemStorage {
	m := MemStorage{
		gaugeMetrics: make(map[string]metric.Gauge, metricsCount),
	}

	return &m
}

func (m *MemStorage) UpdateGauge(name string, val float64) error {

	m.gaugeMetrics[name] = metric.Gauge(val)
	return nil
}

func (m *MemStorage) UpdateCounter(val int64) {
	// m.mux.Lock()
	m.counterMetric += metric.Counter(val)
	// m.mux.Unlock()
}

func (m *MemStorage) GetGauge(name string) (metric.Gauge, error) {
	if val, ok := m.gaugeMetrics[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) GetCounter() (metric.Counter, error) {
	return m.counterMetric, nil
}

func (m *MemStorage) ToList() []string {
	var list []string

	for _, c := range m.gaugeMetrics {
		list = append(list, c.ToString())
	}
	list = append(list, m.counterMetric.ToString())
	return list
}
