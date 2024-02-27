package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/db"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/file"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
)

type StorageCloser interface {
	Save() error
}

type Server struct {
	// можно было бы вообще без этого интерфейса
	// но тогда не понятно - как сохранять метрики в файл в `Run`
	storage   StorageCloser
	webserver *http.Server
	options   *ServerOptions
}

func NewServer(metricsCount int, opt *ServerOptions) *Server {
	var targetStorage handlers.Storage
	var saverStorage StorageCloser
	dbCtx := context.Background()

	if opt.DBDSN == "" {
		if opt.FileStoragePath == "" {
			ts := memory.NewMemStorage(metricsCount)
			saverStorage = ts
			targetStorage = ts
		} else {
			ts := file.NewFileStorage(opt.FileStoragePath, metricsCount, opt.Restore, opt.StoreInterval)
			saverStorage = ts
			targetStorage = ts
		}
	} else {
		ts, err := db.NewDBPool(dbCtx, opt.DBDSN)
		if err != nil {
			logger.Log.Fatal("Не удалось подключиться к БД!")
		}
		targetStorage = ts
		saverStorage = ts
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
				s.webserver.Shutdown(ctx)
				logger.Log.Info("http сервер остановлен!")
				// используется для сохранения метрик в файл
				// но реализован только для файлового стораджа
				// в остальных - методы-заглушки
				err := s.storage.Save()
				if err != nil {
					logger.Log.Error(err)
				}
				return
			case <-ticker.C:
				s.storage.Save()
			}
		}
	}()
	s.webserver.ListenAndServe()
	wg.Wait()
	return nil
}
