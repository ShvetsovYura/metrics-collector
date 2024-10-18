package interceptors

import (
	"context"

	"github.com/ShvetsovYura/metrics-collector/internal/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TrustedSubnetInterceptorWrapper(trustedSubnet string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if trustedSubnet == "" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "в метаданных запроса не найдено значений. x-real-ip обязателен")
		}

		values := md.Get("x-real-ip")
		if len(values) < 1 {
			return nil, status.Error(codes.InvalidArgument, "не укзаан x-real-ip в метаданных запроса")
		}

		xRealIP := values[0]
		if xRealIP == "" {
			return nil, status.Error(codes.InvalidArgument, "клиент находистя вне доверенной подсети")
		}

		val, err := validator.IsIPInSubnet(xRealIP, trustedSubnet)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ошибка при прасинге подсети %s", err)
		}
		if !val {
			return nil, status.Error(codes.PermissionDenied, "доступ запрещен")
		}

		return handler(ctx, req)

	}
}
