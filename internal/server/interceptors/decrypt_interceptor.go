package interceptors

import (
	"context"
	"encoding/json"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DecryptMessage: предназначена для расшифровки входящего тела запроса
func DecryptInterceptor(privateKeyPath string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, inro *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if privateKeyPath != "" {
			request, err := json.Marshal(req)
			if err != nil {
				return nil, status.Error(codes.Internal, "error on read")
			}
			decrytedMessage, err := util.DecryptData(request, privateKeyPath)
			if err != nil {
				return nil, status.Error(codes.Internal, "error on decrypt")
			}
			return handler(ctx, decrytedMessage)
		}
		return handler(ctx, req)
	}
}
