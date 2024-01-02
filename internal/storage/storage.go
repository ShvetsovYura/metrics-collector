package storage

import (
	"fmt"
	"strconv"

	"github.com/ShvetsovYura/metrics-collector/internal/types"
)

type Gauge float64
type Counter int64

func (g Gauge) ToString() string {
	return fmt.Sprintf("%f", g)
}

func (c Counter) ToString() string {
	return fmt.Sprintf("%d", c)
}

type MemStorage struct {
	metrics map[string]types.Stringer
}

func (m *MemStorage) UpdateGauge(name string, val string) {
	v, err := strconv.ParseFloat(val, 64)
	if err == nil {
		m.metrics[name] = Gauge(v)
	}
}

func (m *MemStorage) UpdateCounter(name string, val string) {
	v, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		m.metrics[name] = Counter(v)
	}
}

func (m *MemStorage) GetVal(name string) (types.Stringer, error) {
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
