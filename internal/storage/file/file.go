package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func NewFileStorage(pathToFile string, metricsCount int, restore bool, storeInterval int) *FileStorage {

	immediatelySave := false
	if storeInterval == 0 {
		immediatelySave = true
	}
	s := &FileStorage{
		path:        pathToFile,
		immediately: immediatelySave,
		memStorage:  memory.NewMemStorage(metricsCount),
	}

	if restore {
		s.Restore(context.Background())
	}
	return s
}

func (fs *FileStorage) GetGauge(ctx context.Context, name string) (metric.Gauge, error) {
	return fs.memStorage.GetGauge(ctx, name)
}

func (fs *FileStorage) SetGauge(ctx context.Context, name string, value float64) error {
	if err := fs.memStorage.SetGauge(ctx, name, value); err != nil {
		return err
	}
	fs.SaveNow()
	return nil
}

func (fs *FileStorage) SetCounter(ctx context.Context, name string, value int64) error {
	if err := fs.memStorage.SetCounter(ctx, name, value); err != nil {
		return err
	}
	fs.SaveNow()
	return nil
}

func (fs *FileStorage) GetCounter(ctx context.Context, name string) (metric.Counter, error) {
	return fs.memStorage.GetCounter(ctx, name)
}

func (fs *FileStorage) ToList(ctx context.Context) ([]string, error) {
	return fs.memStorage.ToList(ctx)
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

func (fs *FileStorage) RestoreNow() (map[string]float64, map[string]int64, error) {
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
	logger.Log.Info(di)
	return di.Gauges, di.Counters, nil

}

func (fs *FileStorage) SaveNow() {
	if fs.immediately {
		fs.Save()
	}
}

func (fs *FileStorage) Save() error {
	logger.Log.Info("Начало сохранения метрик в файл ...")
	ctx := context.Background()
	var g = make(map[string]float64, len(fs.memStorage.GetGauges(ctx)))
	var c = make(map[string]int64, len(fs.memStorage.GetCounters(ctx)))
	for k, v := range fs.memStorage.GetGauges(ctx) {
		g[k] = *v.GetRawValue()
	}
	for k, v := range fs.memStorage.GetCounters(ctx) {
		c[k] = *v.GetRawValue()
	}

	err := fs.Dump(g, c)

	if err != nil {
		return err
	}
	logger.Log.Info("Значения метрик успешно сохранены в файл")
	return nil
}

func (fs *FileStorage) Restore(ctx context.Context) error {
	g, c, err := fs.RestoreNow()
	if err != nil {
		return err
	}

	for k, v := range g {
		fs.memStorage.SetGauge(ctx, k, v)
	}
	for k, v := range c {
		fs.memStorage.SetCounter(ctx, k, v)
	}

	return nil
}

func (fs *FileStorage) Ping(ctx context.Context) error {
	return errors.New("it's not db. filestorage")
}

func (fs *FileStorage) SaveGaugesBatch(ctx context.Context, gauges map[string]metric.Gauge) error {
	logger.Log.Info("save metrics in FILE GAUGES")
	return nil
}

func (fs *FileStorage) SaveCountersBatch(ctx context.Context, counters map[string]metric.Counter) error {
	logger.Log.Info("save metrics in FILE COUNTERS")
	return nil
}
