package memory

import (
	"errors"
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
)

type Metric interface {
	ToString() string
}

type MemStorage struct {
	gaugeMetrics  map[string]metric.Gauge
	counterMetric map[string]metric.Counter
}

func NewMemStorage(metricsCount int) *MemStorage {
	m := MemStorage{
		gaugeMetrics:  make(map[string]metric.Gauge, metricsCount),
		counterMetric: make(map[string]metric.Counter, 1),
	}
	return &m
}

func (m *MemStorage) SetGauge(name string, val float64) error {
	m.gaugeMetrics[name] = metric.Gauge(val)
	return nil
}

func (m *MemStorage) SetCounter(name string, val int64) error {
	m.counterMetric[name] += metric.Counter(val)
	return nil
}

func (m *MemStorage) GetGauge(name string) (metric.Gauge, error) {
	if val, ok := m.gaugeMetrics[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) GetCounter(name string) (metric.Counter, error) {
	if val, ok := m.counterMetric[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) GetGauges() map[string]metric.Gauge {
	return m.gaugeMetrics
}

func (m *MemStorage) GetCounters() map[string]metric.Counter {
	return m.counterMetric
}

func (m *MemStorage) ToList() []string {
	var list []string

	for _, c := range m.gaugeMetrics {
		list = append(list, c.ToString())
	}
	for _, c := range m.counterMetric {
		list = append(list, c.ToString())
	}
	return list
}

func (m *MemStorage) Ping() error {
	return errors.New("it's not db. memorystorage")
}

func (m *MemStorage) Save() error {
	return nil
}
func (m *MemStorage) Restore() error {
	return nil
}

func (m *MemStorage) SaveGaugesBatch(gauges map[string]metric.Gauge) {
	return
}
func (m *MemStorage) SaveCountersBatch(couters map[string]metric.Counter) {
	return
}
