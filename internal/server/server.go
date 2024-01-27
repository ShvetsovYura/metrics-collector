package server

import (
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/file"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
)

type Server struct {
	metrics     *memory.MemStorage
	fileStorage *file.FileStorage
	options     *ServerOptions
}

func NewServer(metricsCount int, opt *ServerOptions) *Server {
	fileStorage := file.NewFileStorage("hoho.txt")
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
	router := handlers.ServerRouter(s.metrics)
	return http.ListenAndServe(s.options.EndpointAddr, router)
}
