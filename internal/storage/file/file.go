package file

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
)

type DumpItem struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

type FileStorage struct {
	path        string
	immediately bool
	memStorage  *memory.MemStorage
}

func NewFileStorage(pathToFile string, metricsCount int, immediatelySave bool) *FileStorage {
	return &FileStorage{
		path:        pathToFile,
		immediately: immediatelySave,
		memStorage:  memory.NewMemStorage(metricsCount),
	}
}

func (fs *FileStorage) GetGauge(name string) (metric.Gauge, error) {
	return fs.memStorage.GetGauge(name)
}

func (fs *FileStorage) SetGauge(name string, value float64) error {
	if err := fs.memStorage.SetGauge(name, value); err != nil {
		return err
	}
	fs.SaveNow()
	return nil
}

func (fs *FileStorage) SetCounter(name string, value int64) error {
	if err := fs.memStorage.SetCounter(name, value); err != nil {
		return err
	}
	fs.SaveNow()
	return nil
}

func (fs *FileStorage) GetCounter(name string) (metric.Counter, error) {
	return fs.memStorage.GetCounter(name)
}

func (fs *FileStorage) ToList() []string {
	return fs.memStorage.ToList()
}

func (fs *FileStorage) Dump(gauges map[string]float64, counters map[string]int64) error {
	f, err := os.OpenFile(fs.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	di := DumpItem{Gauges: gauges, Counters: counters}

	data, err := json.MarshalIndent(di, "", "  ")
	if err != nil {
		return err
	}

	f.Write(data)
	return nil
}

func (fs *FileStorage) Restore() (map[string]float64, map[string]int64, error) {
	var buf bytes.Buffer

	f, err := os.OpenFile(fs.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, err
	}

	_, err1 := buf.ReadFrom(f)
	if err1 != nil {
		return nil, nil, err1
	}

	di := DumpItem{}

	err2 := json.Unmarshal(buf.Bytes(), &di)
	if err2 != nil {
		return nil, nil, err2
	}

	return di.Gauges, di.Counters, nil

}

func (fs *FileStorage) SaveNow() {
	if fs.immediately {
		fs.SaveToFile()
	}
}

func (fs *FileStorage) SaveToFile() error {
	logger.Log.Info("Начало сохранения метрик в файл ...")
	var g = make(map[string]float64, len(fs.memStorage.GetGauges()))
	var c = make(map[string]int64, len(fs.memStorage.GetCounters()))
	for k, v := range fs.memStorage.GetGauges() {
		g[k] = *v.GetRawValue()
	}
	for k, v := range fs.memStorage.GetCounters() {
		c[k] = *v.GetRawValue()
	}

	err := fs.Dump(g, c)
	if err != nil {
		return err
	}
	logger.Log.Info("Значения метрик успешно сохранены в файл")
	return nil

}

func (fs *FileStorage) RestoreFromFile() error {
	g, c, err := fs.Restore()
	if err != nil {
		return err
	}

	for k, v := range g {
		fs.memStorage.SetGauge(k, v)
	}
	for k, v := range c {
		fs.memStorage.SetCounter(k, v)
	}

	return nil
}