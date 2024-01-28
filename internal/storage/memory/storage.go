package memory

import (
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/file"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
)

type Metric interface {
	ToString() string
}

type Saver interface {
	Save(map[string]float64, int64) error
}

type MemStorage struct {
	gaugeMetrics    map[string]metric.Gauge
	counterMetric   map[string]metric.Counter
	fs              *file.FileStorage
	immediatelySave bool
}

func NewStorage(metricsCount int, fs *file.FileStorage, immediately bool) *MemStorage {
	m := MemStorage{
		gaugeMetrics:    make(map[string]metric.Gauge, metricsCount),
		counterMetric:   make(map[string]metric.Counter, 1),
		fs:              fs,
		immediatelySave: immediately,
	}
	return &m
}

func (m *MemStorage) UpdateGauge(name string, val float64) error {
	m.gaugeMetrics[name] = metric.Gauge(val)
	return nil
}

func (m *MemStorage) UpdateCounter(name string, val int64) error {
	m.counterMetric[name] += metric.Counter(val)
	logger.Log.Infof("New counter val %v\n", m.counterMetric)
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

func (m *MemStorage) ToList() []string {
	var list []string

	for _, c := range m.gaugeMetrics {
		list = append(list, c.ToString())
	}
	for _, c := range m.counterMetric {
		list = append(list, c.ToString())
	}
	// list = append(list, m.counterMetric.ToString())
	return list
}

func (m *MemStorage) SaveData(s Saver) error {
	var gaugeMetrics map[string]float64 = make(map[string]float64, len(m.gaugeMetrics))
	var counterMetric int64 = int64(m.counterMetric["PollCount"])
	for k, v := range m.gaugeMetrics {
		gaugeMetrics[k] = float64(v)
	}

	return s.Save(gaugeMetrics, counterMetric)
}

func (m *MemStorage) SaveNow() {
	if m.immediatelySave {
		m.SaveToFile()
	}
}

func (m *MemStorage) SaveToFile() error {
	logger.Log.Info("Начало сохранения метрик в файл ...")
	var g = make(map[string]float64, len(m.gaugeMetrics))
	var c = make(map[string]int64, len(m.counterMetric))
	for k, v := range m.gaugeMetrics {
		g[k] = *v.GetRawValue()
	}
	for k, v := range m.counterMetric {
		c[k] = *v.GetRawValue()
	}

	err := m.fs.Dump(g, c)
	if err != nil {
		return err
	}
	logger.Log.Info("Значения метрик успешно сохранены в файл")
	return nil

}

func (m *MemStorage) RestoreFromFile() error {
	g, c, err := m.fs.Restore()
	if err != nil {
		return err
	}

	for k, v := range g {
		m.UpdateGauge(k, v)
	}
	for k, v := range c {
		m.UpdateCounter(k, v)
	}

	return nil
}
