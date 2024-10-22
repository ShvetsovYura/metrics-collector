package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func hashInterceptor(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.New(map[string]string{"HashSHA256": ""})
	outCtx := metadata.NewOutgoingContext(ctx, md)
	errInvoke := invoker(outCtx, method, req, reply, cc, opts...)
	if errInvoke != nil {
		log.Printf("[ERROR] %s,%s", method, errInvoke.Error())
	}
	return errInvoke
}
