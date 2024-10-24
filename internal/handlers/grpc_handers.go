package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/models"
	pb "github.com/ShvetsovYura/metrics-collector/proto"
)

type MetricServer struct {
	pb.UnimplementedMetricsServer
	metrics Storage
}

func NewMetricServer(store Storage) *MetricServer {
	return &MetricServer{metrics: store}
}

// ListMetrics реализует интерфейс получения списка метрик.
func (s *MetricServer) ListMetricsValues(ctx context.Context, in *pb.ListMetricsValuesRequest) (*pb.ListMetricsValuesResponse, error) {
	metricsList, _ := s.metrics.ToList(ctx)
	if metricsList == nil {
		metricsList = []string{}
	}

	return &pb.ListMetricsValuesResponse{
		Values: metricsList,
	}, nil
}

func (s *MetricServer) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse
	logger.Log.Debug("metric type %v", in.Mtype)
	switch in.Mtype {
	case internal.InGaugeName:
		err := s.metrics.SetGauge(ctx, in.Id, float64(in.Value))
		if err != nil {
			logger.Log.Errorf("Ошибка установки значения для gauge: %s, значение: %f. %s", in.Id, in.Value, err.Error())

		}

		currentVal, _ := s.metrics.GetGauge(ctx, in.Id)

		response = pb.UpdateMetricResponse{
			Id:    in.Id,
			Mtype: in.Mtype,
			Value: currentVal.GetRawValue(),
		}

	case internal.InCounterName:
		err := s.metrics.SetCounter(ctx, in.Id, in.Delta)
		if err != nil {
			logger.Log.Errorf("Ошибка установки значения для gauge: %s, значение: %f. %s", in.Id, in.Value, err.Error())
		}

		currentVal, _ := s.metrics.GetCounter(ctx, in.Id)

		response = pb.UpdateMetricResponse{
			Id:    in.Id,
			Mtype: in.Mtype,
			Delta: currentVal.GetRawValue(),
		}

	default:
		logger.Log.Errorf("не удалось определить тип метрики %s", in.Mtype)
	}

	return &response, nil
}

func (s *MetricServer) BatchUpdateMetrics(ctx context.Context, in *pb.BatchUpdateMtericsRequest) (*pb.BatchUpdateMetricsResponse, error) {
	var response pb.BatchUpdateMetricsResponse

	var (
		gauges   = make(map[string]models.Gauge, 100)
		counters = make(map[string]models.Counter, 100)
	)

	for _, mdl := range in.Metrics {
		switch mdl.Mtype {
		case internal.InGaugeName:
			gauges[mdl.Id] = models.Gauge(mdl.Value)
		case internal.InCounterName:
			if v, ok := counters[mdl.Id]; ok {
				counters[mdl.Id] = v + models.Counter(mdl.Delta)
			} else {
				counters[mdl.Id] = models.Counter(mdl.Delta)
			}
		}
	}
	var err error
	err = s.metrics.SaveCountersBatch(ctx, counters)
	if err != nil {
		logger.Log.Error(err.Error())
	}

	err = s.metrics.SaveGaugesBatch(ctx, gauges)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	response = pb.BatchUpdateMetricsResponse{}

	return &response, nil
}

func (s *MetricServer) GetMetric(ctx context.Context, in *pb.GetMetricRequest) (*pb.GetMetricResponse, error) {
	valGauge, err := s.metrics.GetGauge(ctx, in.Name)
	if err == nil {
		return &pb.GetMetricResponse{
			Id:    in.Name,
			Mtype: "gauge",
			Value: valGauge.GetRawValue(),
		}, nil
	} else {
		logger.Log.Error(err.Error())
	}

	valCounter, err := s.metrics.GetCounter(ctx, in.Name)
	if err != nil {
		return nil, errors.New("не найдена мертика по такому имени")
	}
	return &pb.GetMetricResponse{
		Id:    in.Name,
		Mtype: "counter",
		Delta: valCounter.GetRawValue(),
	}, nil

}

func (s *MetricServer) DBPing(ctx context.Context, in *pb.DbPingRequest) (*pb.DbPingResponse, error) {
	err := s.metrics.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("бд недоступна %w", err)
	}
	return &pb.DbPingResponse{}, nil
}
