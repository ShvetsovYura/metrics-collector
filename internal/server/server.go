package server

import (
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/storage/memory"
)

type Server struct {
	metricCount int
	options     *ServerOptions
}

func NewServer(opt *ServerOptions) *Server {
	return &Server{
		metricCount: 40,
		options:     opt,
	}
}

func (s *Server) Run() error {
	m := memory.NewStorage(s.metricCount)
	router := handlers.ServerRouter(m)
	return http.ListenAndServe(s.options.EndpointAddr, router)
}
