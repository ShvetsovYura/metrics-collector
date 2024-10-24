package grpcclient

import (
	"context"
	"fmt"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
	pb "github.com/ShvetsovYura/metrics-collector/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type GRPCClient struct {
	conn    *grpc.ClientConn
	client  pb.MetricsClient
	hashKey string
}

func NewClient(addr string, hashKey string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать GRPC клиент %w", err)
	}
	return &GRPCClient{
		conn:    conn,
		client:  pb.NewMetricsClient(conn),
		hashKey: hashKey,
	}, nil
}

func (g *GRPCClient) Close() {
	g.conn.Close()
}

func (g *GRPCClient) Send(item agent.MetricItem) error {
	var respHeaders metadata.MD
	logger.Log.Debug("grpc start send metric")
	msg := pb.UpdateMetricRequest{
		Id:    item.ID,
		Mtype: item.MType,
		Value: item.Value,
		Delta: item.Delta,
	}

	md := metadata.New(map[string]string{})
	if g.hashKey != "" {
		msgData, err := proto.Marshal(&msg)
		if err != nil {
			return err
		}
		md.Append("HashSHA256", util.Hash(msgData, g.hashKey))
	}
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	logger.Log.Debug("before send")
	resp, err := g.client.UpdateMetric(ctx, &msg,
		grpc.Header(&respHeaders), grpc.UseCompressor(gzip.Name))
	logger.Log.Debug("after send %v", resp)
	if err != nil {
		return fmt.Errorf("не удалось отправить метрики, %w", err)
	}
	logger.Log.Debug("end send metric")
	return nil
}
