package storage

import (
	"fmt"
	"strconv"
)

type Gauge float64
type Counter int64

type Metric interface {
	ToString() string
}

func (g Gauge) ToString() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func (c Counter) ToString() string {
	return strconv.FormatInt(int64(c), 10)
}

type MemStorage struct {
	metrics map[string]Metric
}

func NewStorage(metricsCount int) *MemStorage {
	ms := new(MemStorage)
	ms.metrics = make(map[string]Metric, metricsCount)
	return ms
}

func (m *MemStorage) UpdateGauge(name string, val float64) {
	m.metrics[name] = Gauge(val)
}

func (m *MemStorage) UpdateCounter(name string, val int64) {
	if v, ok := m.metrics[name]; ok {
		m.metrics[name] = v.(Counter) + Counter(val)
	} else {
		m.metrics[name] = Counter(val)
	}
}

func (m *MemStorage) GetVal(name string) (Metric, error) {
	if val, ok := m.metrics[name]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("NotFound %s", name)
	}
}

func (m *MemStorage) ToList() []string {
	var list []string

	for _, c := range m.metrics {
		list = append(list, c.ToString())
	}
	return list
}
