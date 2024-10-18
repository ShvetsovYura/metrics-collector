package interceptors

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func hasInterceptor(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(map[string]string{"HashSHA256": ""})
	outCtx := metadata.NewOutgoingContext(ctx, md)
	errInvoke := invoker(outCtx, method, req, reply, cc, opts...)
	// выполняем действия после вызова метода
	if errInvoke != nil {
		log.Printf("[ERROR] %s,%s", method, errInvoke.Error())
	}
	return errInvoke
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
