package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

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

func (m *MemStorage) SetGauge(ctx context.Context, name string, val float64) error {
	m.mx.Lock()
	m.gaugeMetrics[name] = metric.Gauge(val)
	m.mx.Unlock()
	return nil
}

func (m *MemStorage) SetGauges(ctx context.Context, gauges map[string]float64) {
	for k, v := range gauges {
		m.SetGauge(ctx, k, v)
	}
}

func (m *MemStorage) SetCounters(ctx context.Context, counters map[string]int64) {
	for k, v := range counters {
		m.SetCounter(ctx, k, v)
	}
}

func (m *MemStorage) SetCounter(ctx context.Context, name string, val int64) error {
	m.mx.Lock()
	m.counterMetric[name] += metric.Counter(val)
	m.mx.Unlock()
	return nil
}

func (m *MemStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, error) {
	m.mx.Lock()
	val, ok := m.gaugeMetrics[name]
	m.mx.Unlock()
	if ok {
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

func (m *MemStorage) Ping(ctx context.Context) error {
	return nil
	// return errors.New("it's not db. memorystorage")
}

func (m *MemStorage) SaveGaugesBatch(ctx context.Context, gauges map[string]metric.Gauge) error {
	return nil
}
func (m *MemStorage) SaveCountersBatch(ctx context.Context, couters map[string]metric.Counter) error {
	return nil
}
func (m *MemStorage) Save() error {
	return nil
}
func (m *MemStorage) Restore(ctx context.Context) error {
	return nil
}
