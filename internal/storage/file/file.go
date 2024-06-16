package file

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/metric"
)

type MemoryStore interface {
	GetGauge(ctx context.Context, name string) (metric.Gauge, error)
	GetGauges(ctx context.Context) map[string]metric.Gauge
	SetGauge(ctx context.Context, name string, val float64) error
	SetGauges(ctx context.Context, gauges map[string]float64)
	GetCounter(ctx context.Context, name string) (metric.Counter, error)
	GetCounters(ctx context.Context) map[string]metric.Counter
	SetCounter(ctx context.Context, name string, value int64) error
	SetCounters(ctx context.Context, gauges map[string]int64)

	ToList(ctx context.Context) ([]string, error)
}

type DumpItem struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

type FileStorage struct {
	path        string
	immediately bool
	memStorage  MemoryStore
}

func NewFileStorage(pathToFile string, memStorage MemoryStore, restore bool, storeInterval int) *FileStorage {

	immediatelySave := false
	if storeInterval == 0 {
		immediatelySave = true
	}
	s := &FileStorage{
		path:        pathToFile,
		immediately: immediatelySave,
		memStorage:  memStorage,
	}

	if restore {
		err := s.Restore(context.Background())
		if err != nil {
			logger.Log.Errorf("Ошибка при восстановлении метрик из файла, %s", err.Error())
		}
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
	defer func() {
		err := f.Close()
		if err != nil {
			logger.Log.Errorf("ошибка закрытия файла, %s", err.Error())
		}
	}()
	di := DumpItem{Gauges: gauges, Counters: counters}

	data, err := json.MarshalIndent(di, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл, %w", err)
	}
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
		err := fs.Save()
		if err != nil {
			logger.Log.Errorf("ошибка сохранения метрик в файл, %s", err.Error())
		}
	}
}

func (fs *FileStorage) ExtractGauges(ctx context.Context) map[string]float64 {
	gauges := fs.memStorage.GetGauges(ctx)
	var g = make(map[string]float64, len(gauges))
	for k, v := range gauges {
		g[k] = *v.GetRawValue()
	}
	return g
}

func (fs *FileStorage) ExtractCounters(ctx context.Context) map[string]int64 {
	counters := fs.memStorage.GetCounters(ctx)
	var c = make(map[string]int64, len(counters))
	for k, v := range counters {
		c[k] = *v.GetRawValue()
	}
	return c
}

func (fs *FileStorage) Save() error {
	logger.Log.Info("Начало сохранения метрик в файл ...")
	ctx := context.Background()
	g := fs.ExtractGauges(ctx)
	c := fs.ExtractCounters(ctx)
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

	fs.memStorage.SetGauges(ctx, g)
	fs.memStorage.SetCounters(ctx, c)

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
