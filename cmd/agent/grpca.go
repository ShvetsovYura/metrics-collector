package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/ShvetsovYura/metrics-collector/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc/metadata"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	var respHeaders metadata.MD
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewMetricsClient(conn)
	md := metadata.New(map[string]string{"HashSHA256": "066985110483cecc7b9e52576c2852829a3886c1eeff6dfe5cd94034805f307a"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	r, err := c.GetMetric(ctx, &pb.GetMetricRequest{
		Name: "hoho",
	}, grpc.Header(&respHeaders), grpc.UseCompressor(gzip.Name))

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(r)
}
