package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
)

type Metric interface {
	ToString() string
}

type MemStorage struct {
	mx            sync.Mutex
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

func (m *MemStorage) SetGauge(_ context.Context, name string, val float64) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.gaugeMetrics[name] = metric.Gauge(val)
	return nil
}

func (m *MemStorage) SetGauges(_ context.Context, gauges map[string]float64) {
	for k, v := range gauges {
		err := m.SetGauge(context.TODO(), k, v)
		if err != nil {
			logger.Log.Errorf("ошибка при записи gauge: %s:%d", k, v)
		}
	}
}

func (m *MemStorage) SetCounters(_ context.Context, counters map[string]int64) {
	for k, v := range counters {
		err := m.SetCounter(context.TODO(), k, v)
		if err != nil {
			logger.Log.Errorf("ошибка при записи counter: %s:%d", k, v)
		}
	}
}

func (m *MemStorage) SetCounter(_ context.Context, name string, val int64) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.counterMetric[name] += metric.Counter(val)
	return nil
}

func (m *MemStorage) GetGauge(_ context.Context, name string) (metric.Gauge, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if val, ok := m.gaugeMetrics[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) GetCounter(_ context.Context, name string) (metric.Counter, error) {
	m.mx.Lock()
	defer m.mx.Unlock()
	if val, ok := m.counterMetric[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) GetGauges(_ context.Context) map[string]metric.Gauge {
	return m.gaugeMetrics
}

func (m *MemStorage) GetCounters(_ context.Context) map[string]metric.Counter {
	return m.counterMetric
}

func (m *MemStorage) ToList(_ context.Context) ([]string, error) {
	var list []string

	gaugeKeys := make([]string, 0, len(m.gaugeMetrics))
	counterKeys := make([]string, 0, len(m.counterMetric))
	for key := range m.gaugeMetrics {
		gaugeKeys = append(gaugeKeys, key)
	}
	for key := range m.counterMetric {
		counterKeys = append(counterKeys, key)
	}
	sort.Strings(gaugeKeys)
	sort.Strings(counterKeys)

	for _, k := range gaugeKeys {
		if v, ok := m.gaugeMetrics[k]; ok {
			list = append(list, v.ToString())
		}
	}
	for _, k := range counterKeys {
		if v, ok := m.counterMetric[k]; ok {
			list = append(list, v.ToString())
		}
	}
	return list, nil
}

func (m *MemStorage) Ping(_ context.Context) error {
	return nil
}

func (m *MemStorage) SaveGaugesBatch(_ context.Context, gauges map[string]metric.Gauge) error {
	return nil
}
func (m *MemStorage) SaveCountersBatch(_ context.Context, couters map[string]metric.Counter) error {
	return nil
}
func (m *MemStorage) Save() error {
	return nil
}
func (m *MemStorage) Restore(_ context.Context) error {
	return nil
}
