package interceptors

import (
	"compress/gzip"
	"fmt"
	"io"
)

type ComporessReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*ComporessReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации ридера, %w", err)
	}

	return &ComporessReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *ComporessReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)

}

func (c *ComporessReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия оригинального ридера, %w", err)
	}
	if err := c.zr.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия gzip ридера, %w", err)
	}
	return nil
}

// func Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if ok {
// 		values := md.Get("encoding")
// 		if len(values) > 0 && slices.Contains(values, "zip") {
// 			cr, err := NewCompressReader(req.(io.ReadCloser))
// 		}
// 	}

// }
