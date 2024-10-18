package grpcclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ShvetsovYura/metrics-collector/internal/agent"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn      *grpc.ClientConn
	hashKey   string
	cryptoKey string
}

func NewClient(addr string, hashKey string, cryptoKey string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		encryptData,
	))
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать GRPC клиент %w", err)
	}
	return &GRPCClient{
		conn:      conn,
		hashKey:   hashKey,
		cryptoKey: cryptoKey,
	}, nil
}

func (g *GRPCClient) Close() {
	g.conn.Close()
}

func (g *GRPCClient) Send(item agent.MetricItem) error {

	// var respHeaders metadata.MD
	// // получаем переменную интерфейсного типа UsersClient,
	// // через которую будем отправлять сообщения
	// c := pb.NewMetricsClient(conn)
	// md := metadata.New(map[string]string{"HashSHA256": "066985110483cecc7b9e52576c2852829a3886c1eeff6dfe5cd94034805f307a"})
	// ctx := metadata.NewOutgoingContext(context.Background(), md)

	// r, err := c.GetMetric(ctx, &pb.GetMetricRequest{
	// 	Name: "hoho",
	// }, grpc.Header(&respHeaders), grpc.UseCompressor(gzip.Name))

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	fmt.Println("send")
	return nil
}
func main() {
	// устанавливаем соединение с сервером

}

func encryptData(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	d, _ := os.Getwd()
	p := filepath.Join(d, "public.pem")
	jsonData, _ := json.Marshal(req)
	data, err := util.EncryptData(jsonData, p)

	errInvoke := invoker(ctx, method, data, reply, cc, opts...)

	// выполняем действия после вызова метода
	if errInvoke != nil {
		log.Printf("[ERROR] %s,%s", method, errInvoke.Error())
	}
	return err
}
