package server

import (
	"context"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/db"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/file"
)

type Server struct {
	storage *file.FileStorage
	dbPool  *db.DB
	options *ServerOptions
}

func NewServer(metricsCount int, opt *ServerOptions) *Server {
	immediatelySave := false
	if opt.StoreInterval == 0 {
		immediatelySave = true
	}
	fileStorage := file.NewFileStorage(opt.FileStoragePath, metricsCount, immediatelySave)
	dbCtx := context.Background()
	dbStorage, err := db.NewDBPool(dbCtx, opt.DbDSN)
	if err != nil {
		panic(err)
	}
	return &Server{
		storage: fileStorage,
		dbPool:  dbStorage,
		options: opt,
	}
}

func (s *Server) Run() error {
	if s.options.Restore {
		s.storage.RestoreFromFile()
	}

	router := handlers.ServerRouter(s.storage, s.dbPool)
	srv := &http.Server{
		Addr:    s.options.EndpointAddr,
		Handler: router,
	}
	ticker := time.NewTicker(time.Duration(s.options.StoreInterval) * time.Second)
	// https://habr.com/ru/articles/771626/
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

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
				s.storage.SaveToFile()
				return
			case <-ticker.C:
				s.storage.SaveToFile()
			}
		}
	}()

	srv.ListenAndServe()
	wg.Wait()
	return nil
}
