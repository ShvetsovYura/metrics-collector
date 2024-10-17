package server

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/handlers"
	"github.com/ShvetsovYura/metrics-collector/internal/interceptors"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/storage"
	pb "github.com/ShvetsovYura/metrics-collector/proto"
	"google.golang.org/grpc"
)

type StorageCloser interface {
	Save() error
}

type IServer interface {
	StartListen() error
	Shutdown(ctx context.Context) error
	RegisterHandlers(handlers.Storage, *Options)
}

// Server, хранит информации о сервере сбора метрик.
type Server struct {
	// можно было бы вообще без этого интерфейса
	// но тогда не понятно - как сохранять метрики в файл в `Run`
	storage StorageCloser
	server  IServer
	options *Options
}

// NewServer, создает новый сервер работы с метриками.
func NewServer(server IServer, metricsCount int, opt *Options) *Server {
	dbCtx := context.Background()

	var (
		targetStorage handlers.Storage
		saverStorage  StorageCloser
	)
	// TODO: Подумать над упрощением
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
	server.RegisterHandlers(targetStorage, opt)
	return &Server{
		// из-за того, что удалил методы Save и Restore из интерфейса Storage
		// приходится костылить такое - дублирование стораджа, но с другим интерфейсом
		storage: saverStorage,
		server:  server,
		options: opt,
	}
}

// Run, запускает сервер.
func (s *Server) Run(ctx context.Context) error {
	logger.Log.Info("run Server app")
	var wg sync.WaitGroup
	wg.Add(1)
	ticker := time.NewTicker(s.options.StoreInterval)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Останавливаю сервер...")
				// сначала - остановка сервера, чтобы не принимал новые запросы
				if err := s.server.Shutdown(ctx); err != nil {
					logger.Log.Fatalf("не удалось остановить сервер %s", err.Error())
				}

				// используется для сохранения метрик в файл
				// но реализован только для файлового стораджа
				// в остальных - методы-заглушки
				if err := s.storage.Save(); err != nil {
					logger.Log.Error(err)
				}
				logger.Log.Info("http сервер остановлен!")
				return
			case <-ticker.C:
				logger.Log.Debug("что-то другое")
				err := s.storage.Save()
				if err != nil {
					logger.Log.Errorf("Ошибка сохранения метрик, %s", err.Error())
				}
			}
		}
	}()

	if err := s.server.StartListen(); err != http.ErrServerClosed {
		logger.Log.Fatalf("не удалось запусить web сервер, %s", err.Error())
	}

	wg.Wait()
	return nil
}

type HTTPServer struct {
	webserver *http.Server
}

func NewHttpServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) StartListen() error {
	err := s.webserver.ListenAndServe()
	return err
}

func (s *HTTPServer) RegisterHandlers(targetStorage handlers.Storage, opt *Options) {
	s.webserver = &http.Server{
		Addr:    opt.EndpointAddr,
		Handler: handlers.ServerRouter(targetStorage, opt.Key, opt.CryptoKey, opt.TrustedSubnet),
	}
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	err := s.webserver.Shutdown(ctx)
	return err
}

type GRPCServer struct {
	grpcServer grpc.Server
}

func NewGRPCServer() *GRPCServer {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.HashInterceptorWrapper("abracadabra"),
		),
	}
	return &GRPCServer{
		grpcServer: *grpc.NewServer(opts...),
	}
}

func (s *GRPCServer) RegisterHandlers(targetStorage handlers.Storage, opt *Options) {
	pb.RegisterMetricsServer(
		&s.grpcServer,
		handlers.NewMetricServer(targetStorage),
	)
}

func (s *GRPCServer) StartListen() error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(listen)

}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	s.grpcServer.Stop()
	return nil
}
