package storage

import (
	"fmt"
	"strconv"

	"github.com/ShvetsovYura/metrics-collector/internal/types"
)

type Gauge float64
type Counter int64

func (g Gauge) ToString() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func (c Counter) ToString() string {
	return fmt.Sprintf("%d", c)
}

type MemStorage struct {
	metrics map[string]types.Stringer
}

func New() MemStorage {
	ms := MemStorage{}
	ms.metrics = make(map[string]types.Stringer, 40)
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
