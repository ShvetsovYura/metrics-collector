package memory

import (
	"context"
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

func (m *MemStorage) SetGauge(ctx context.Context, name string, val float64) error {
	m.gaugeMetrics[name] = metric.Gauge(val)
	return nil
}

func (m *MemStorage) SetCounter(ctx context.Context, name string, val int64) error {
	m.counterMetric[name] += metric.Counter(val)
	return nil
}

func (m *MemStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, error) {
	if val, ok := m.gaugeMetrics[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) GetCounter(ctx context.Context, name string) (metric.Counter, error) {
	if val, ok := m.counterMetric[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) GetGauges(ctx context.Context) map[string]metric.Gauge {
	return m.gaugeMetrics
}

func (m *MemStorage) GetCounters(ctx context.Context) map[string]metric.Counter {
	return m.counterMetric
}

func (m *MemStorage) ToList(ctx context.Context) ([]string, error) {
	var list []string

	for _, c := range m.gaugeMetrics {
		list = append(list, c.ToString())
	}
	for _, c := range m.counterMetric {
		list = append(list, c.ToString())
	}
	return list, nil
}

func (m *MemStorage) Ping(ctx context.Context) error {
	return errors.New("it's not db. memorystorage")
}

func (m *MemStorage) Save() error {
	return nil
}
func (m *MemStorage) Restore(ctx context.Context) error {
	return nil
}

func (m *MemStorage) SaveGaugesBatch(ctx context.Context, gauges map[string]metric.Gauge) error {
	return nil
}
func (m *MemStorage) SaveCountersBatch(ctx context.Context, couters map[string]metric.Counter) error {
	return nil
}
