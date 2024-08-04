package middlewares

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации ридера, %w", err)
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
	// if err!=nil{
	// 	return num, fmt.Errorf("gzip read error %w", err)
	// }
	// }
	// return num, nil
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия оригинального ридера, %w", err)
	}
	if err := c.zr.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия gzip ридера, %w", err)
	}
	return nil
}

// WithUnzipRequest, мидлваря для распаковки принятых сжатых данных.
func WithUnzipRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// распаковка входящих сжатых данных
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr

			defer func() {
				err := cr.Close()
				if err != nil {
					logger.Log.Errorf("ошибка закрытия reader, %s", err.Error())
				}
			}()
		}

		next.ServeHTTP(w, r)
	})
}
