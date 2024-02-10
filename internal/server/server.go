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
)

type Server struct {
	storage handlers.Storage
	options *ServerOptions
}

func NewServer(metricsCount int, opt *ServerOptions) *Server {

	var targetStorage handlers.Storage
	fileStorage := file.NewFileStorage(opt.FileStoragePath, metricsCount, opt.Restore, opt.StoreInterval)
	targetStorage = fileStorage
	dbCtx := context.Background()
	dbStorage, err := db.NewDBPool(dbCtx, opt.DBDSN)
	if err == nil {
		pingErr := dbStorage.Ping()
		if pingErr == nil {
			targetStorage = dbStorage
		}
	}

	return &Server{
		storage: targetStorage,
		options: opt,
	}
}

func (s *Server) Run(ctx context.Context) error {

	srv := &http.Server{
		Addr:    s.options.EndpointAddr,
		Handler: handlers.ServerRouter(s.storage),
	}

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
				srv.Shutdown(ctx)
				logger.Log.Info("http сервер остановлен!")
				s.storage.Save()
				return
			case <-ticker.C:
				s.storage.Save()
			}
		}
	}()

	srv.ListenAndServe()
	wg.Wait()
	return nil
}
