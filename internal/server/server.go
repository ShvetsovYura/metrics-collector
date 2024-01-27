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
	"github.com/ShvetsovYura/metrics-collector/internal/storage/file"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
)

type Server struct {
	metrics     *memory.MemStorage
	fileStorage *file.FileStorage
	options     *ServerOptions
}

func NewServer(metricsCount int, opt *ServerOptions) *Server {
	fileStorage := file.NewFileStorage(opt.FileStoragePath)
	immediatelySave := false
	if opt.StoreInterval == 0 {
		immediatelySave = true
	}
	memStorage := memory.NewStorage(metricsCount, fileStorage, immediatelySave)
	return &Server{
		metrics:     memStorage,
		fileStorage: fileStorage,
		options:     opt,
	}
}

func (s *Server) Run() error {
	if s.options.Restore {
		s.metrics.RestoreFromFile()
	}

	router := handlers.ServerRouter(s.metrics)
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
				s.metrics.SaveToFile()
				return
			case <-ticker.C:
				s.metrics.SaveToFile()
			}
		}
	}()

	srv.ListenAndServe()
	wg.Wait()
	return nil
}
