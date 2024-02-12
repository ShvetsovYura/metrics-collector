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
	storage   StorageCloser
	webserver *http.Server
	options   *ServerOptions
}

func NewServer(metricsCount int, opt *ServerOptions) *Server {

	var targetStorage handlers.Storage

	dbCtx := context.Background()

	if opt.DBDSN == "" {
		if opt.FileStoragePath == "" {
			targetStorage = memory.NewMemStorage(metricsCount)
		} else {
			targetStorage = file.NewFileStorage(opt.FileStoragePath, metricsCount, opt.Restore, opt.StoreInterval)
		}
	} else {
		dbStorage, err := db.NewDBPool(dbCtx, opt.DBDSN)
		if err != nil {
			logger.Log.Fatal("Не удалось подключиться к БД!")
		}
		targetStorage = dbStorage
	}
	return &Server{
		storage: targetStorage,
		webserver: &http.Server{
			Addr:    opt.EndpointAddr,
			Handler: handlers.ServerRouter(targetStorage),
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
