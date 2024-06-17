package storage

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
)

type Metric interface {
	ToString() string
}

type Memory struct {
	mx            sync.Mutex
	gaugeMetrics  map[string]models.Gauge
	counterMetric map[string]models.Counter
}

func NewMemory(metricsCount int) *Memory {
	m := Memory{
		gaugeMetrics:  make(map[string]models.Gauge, metricsCount),
		counterMetric: make(map[string]models.Counter, 1),
	}
	return &m
}

func (m *Memory) SetGauge(_ context.Context, name string, val float64) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.gaugeMetrics[name] = models.Gauge(val)
	return nil
}

func (m *Memory) SetGauges(_ context.Context, gauges map[string]float64) {
	for k, v := range gauges {
		err := m.SetGauge(context.TODO(), k, v)
		if err != nil {
			logger.Log.Errorf("ошибка при записи gauge: %s:%d", k, v)
		}
	}
}

func (m *Memory) SetCounters(_ context.Context, counters map[string]int64) {
	for k, v := range counters {
		err := m.SetCounter(context.TODO(), k, v)
		if err != nil {
			logger.Log.Errorf("ошибка при записи counter: %s:%d", k, v)
		}
	}
}

func (m *Memory) SetCounter(_ context.Context, name string, val int64) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.counterMetric[name] += models.Counter(val)
	return nil
}

func (m *Memory) GetGauge(_ context.Context, name string) (models.Gauge, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if val, ok := m.gaugeMetrics[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *Memory) GetCounter(_ context.Context, name string) (models.Counter, error) {
	m.mx.Lock()
	defer m.mx.Unlock()
	if val, ok := m.counterMetric[name]; ok {
		return val, nil
	} else {
		return 0, fmt.Errorf("NotFound %s", name)
	}
}

func (m *Memory) GetGauges(_ context.Context) map[string]models.Gauge {
	return m.gaugeMetrics
}

func (m *Memory) GetCounters(_ context.Context) map[string]models.Counter {
	return m.counterMetric
}

func (m *Memory) ToList(_ context.Context) ([]string, error) {
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

func (m *Memory) Ping(_ context.Context) error {
	return nil
}

func (m *Memory) SaveGaugesBatch(_ context.Context, gauges map[string]models.Gauge) error {
	return nil
}
func (m *Memory) SaveCountersBatch(_ context.Context, couters map[string]models.Counter) error {
	return nil
}
func (m *Memory) Save() error {
	return nil
}
func (m *Memory) Restore(_ context.Context) error {
	return nil
}
