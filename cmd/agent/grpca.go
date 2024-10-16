package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/ShvetsovYura/metrics-collector/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewMetricsClient(conn)

	r, err := c.SetMetric(context.Background(), &pb.SetMetricRequest{
		Metric: &pb.Metric{
			Id:    "AwesomeMetric",
			Mtype: "gauge",
			Value: float32(123.33),
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(r)
}
