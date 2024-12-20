package interceptors

import (
	"context"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func HashInterceptorWrapper(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			values := md.Get("HashSHA256")
			if len(values) > 0 {
				hashHeader := values[0]
				if key != "" && hashHeader != "" {
					body, _ := proto.Marshal(req.(proto.Message))
					hash := util.Hash(body, key)
					if hashHeader != hash {
						logger.Log.Infof("key %s hashHeader: %s hash: %s", key, hashHeader, hash)
						return nil, status.Error(codes.InvalidArgument, "hashes not equal")
					}
				}
			}
		}
		res, err := handler(ctx, req)
		body, _ := proto.Marshal(res.(proto.Message))
		hash := util.Hash(body, key)
		respMd := metadata.New(map[string]string{"HashSHA256": hash})
		if err := grpc.SendHeader(ctx, respMd); err != nil {
			return nil, status.Error(codes.Internal, "unable to send 'HashSHA256' header")
		}

		return res, err
	}
}
