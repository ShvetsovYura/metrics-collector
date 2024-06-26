package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
)

type StorageCloser interface {
	Save() error
}

// Server, хранит информации о сервере сбора метрик.
type Server struct {
	// можно было бы вообще без этого интерфейса
	// но тогда не понятно - как сохранять метрики в файл в `Run`
	storage   StorageCloser
	webserver *http.Server
	options   *Options
}

// NewServer, создает новый сервер работы с метриками.
func NewServer(metricsCount int, opt *Options) *Server {
	dbCtx := context.Background()

	var (
		targetStorage handlers.Storage
		saverStorage  StorageCloser
	)

	if opt.DBDSN == "" {
		m := storage.NewMemory(metricsCount)
		if opt.FileStoragePath == "" {
			saverStorage = m
			targetStorage = m
		} else {
			f := storage.NewFile(opt.FileStoragePath, m, opt.Restore, opt.StoreInterval)
			saverStorage = f
			targetStorage = f
		}
	} else {
		d, err := storage.NewDBPool(dbCtx, opt.DBDSN)
		if err != nil {
			logger.Log.Fatal("Не удалось подключиться к БД!")
		}

		targetStorage = d
		saverStorage = d
	}

	return &Server{
		// из-за того, что удалил методы Save и Restore из интерфейса Storage
		// приходится костылить такое - дублирование стораджа, но с другим интерфейсом
		storage: saverStorage,
		webserver: &http.Server{
			Addr:    opt.EndpointAddr,
			Handler: handlers.ServerRouter(targetStorage, opt.Key),
		},
		options: opt,
	}
}

// Run, запускает сервер.
func (s *Server) Run(ctx context.Context) error {
	logger.Log.Info("START HTTP SERVER")

	ticker := time.NewTicker(time.Duration(s.options.StoreInterval) * time.Second)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Останавливаю http сервер...")

				err := s.webserver.Shutdown(ctx)
				if err != nil {
					logger.Log.Fatalf("не удалось остановить сервер %s", err.Error())
				}

				logger.Log.Info("http сервер остановлен!")
				// используется для сохранения метрик в файл
				// но реализован только для файлового стораджа
				// в остальных - методы-заглушки
				err = s.storage.Save()
				if err != nil {
					logger.Log.Error(err)
				}

				return
			case <-ticker.C:
				err := s.storage.Save()
				if err != nil {
					logger.Log.Errorf("Ошибка сохранения метрик, %s", err.Error())
				}
			}
		}
	}()

	err := s.webserver.ListenAndServe()
	logger.Log.Fatalf("не удалось запусить web сервер, %s", err.Error())
	wg.Wait()

	return nil
}
