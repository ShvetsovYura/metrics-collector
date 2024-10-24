package interceptors

import (
	"context"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func HashInterceptorWrapper(key string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		data, err := proto.Marshal(req.(proto.Message))
		if err != nil {
			return err
		}
		hash := util.Hash(data, key)
		outCtx := metadata.AppendToOutgoingContext(ctx, "HashSHA256", hash)
		return invoker(outCtx, method, req, reply, cc, opts...)
	}
}
