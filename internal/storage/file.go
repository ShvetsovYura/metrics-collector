package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
)

type MemoryStore interface {
	GetGauge(ctx context.Context, name string) (models.Gauge, error)
	GetGauges(ctx context.Context) map[string]models.Gauge
	SetGauge(ctx context.Context, name string, val float64) error
	SetGauges(ctx context.Context, gauges map[string]float64)
	GetCounter(ctx context.Context, name string) (models.Counter, error)
	GetCounters(ctx context.Context) map[string]models.Counter
	SetCounter(ctx context.Context, name string, value int64) error
	SetCounters(ctx context.Context, gauges map[string]int64)
	ToList(ctx context.Context) ([]string, error)
}

type File struct {
	path        string
	immediately bool
	memStorage  MemoryStore
}

func NewFile(pathToFile string, memStorage MemoryStore, restore bool, storeInterval int) *File {
	immediatelySave := false
	if storeInterval == 0 {
		immediatelySave = true
	}

	s := &File{
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

func (fs *File) GetGauge(ctx context.Context, name string) (models.Gauge, error) {
	val, err := fs.memStorage.GetGauge(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return val, nil
}

func (fs *File) SetGauge(ctx context.Context, name string, value float64) error {
	if err := fs.memStorage.SetGauge(ctx, name, value); err != nil {
		return fmt.Errorf("%w", err)
	}

	fs.SaveNow()

	return nil
}

func (fs *File) SetCounter(ctx context.Context, name string, value int64) error {
	if err := fs.memStorage.SetCounter(ctx, name, value); err != nil {
		return fmt.Errorf("%w", err)
	}

	fs.SaveNow()

	return nil
}

func (fs *File) GetCounter(ctx context.Context, name string) (models.Counter, error) {
	val, err := fs.memStorage.GetCounter(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return val, nil
}

func (fs *File) ToList(ctx context.Context) ([]string, error) {
	val, err := fs.memStorage.ToList(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return val, nil
}

func (fs *File) Dump(gauges map[string]float64, counters map[string]int64) error {
	f, err := os.OpenFile(fs.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			logger.Log.Errorf("ошибка закрытия файла, %s", closeErr.Error())
		}
	}()

	di := models.DumpItem{Gauges: gauges, Counters: counters}

	data, err := json.MarshalIndent(di, "", "  ")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл, %w", err)
	}

	return nil
}

func (fs *File) RestoreNow() (map[string]float64, map[string]int64, error) {
	var buf bytes.Buffer

	f, err := os.OpenFile(fs.path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка открытия файла, %w", err)
	}

	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка чтения файла, %w", err)
	}

	di := models.DumpItem{}

	err = json.Unmarshal(buf.Bytes(), &di)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка преборазования из json, %w", err)
	}

	logger.Log.Info(di)

	return di.Gauges, di.Counters, nil
}

func (fs *File) SaveNow() {
	if fs.immediately {
		err := fs.Save()
		if err != nil {
			logger.Log.Errorf("ошибка сохранения метрик в файл, %s", err.Error())
		}
	}
}

func (fs *File) ExtractGauges(ctx context.Context) map[string]float64 {
	gauges := fs.memStorage.GetGauges(ctx)

	var g = make(map[string]float64, len(gauges))

	for k, v := range gauges {
		g[k] = *v.GetRawValue()
	}

	return g
}

func (fs *File) ExtractCounters(ctx context.Context) map[string]int64 {
	counters := fs.memStorage.GetCounters(ctx)

	var c = make(map[string]int64, len(counters))

	for k, v := range counters {
		c[k] = *v.GetRawValue()
	}

	return c
}

func (fs *File) Save() error {
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

func (fs *File) Restore(ctx context.Context) error {
	g, c, err := fs.RestoreNow()
	if err != nil {
		return err
	}

	fs.memStorage.SetGauges(ctx, g)
	fs.memStorage.SetCounters(ctx, c)

	return nil
}

func (fs *File) Ping(_ context.Context) error {
	return errors.New("it's not db. filestorage")
}

func (fs *File) SaveGaugesBatch(_ context.Context, _ map[string]models.Gauge) error {
	logger.Log.Info("save metrics in FILE GAUGES")
	return nil
}

func (fs *File) SaveCountersBatch(_ context.Context, _ map[string]models.Counter) error {
	logger.Log.Info("save metrics in FILE COUNTERS")
	return nil
}
