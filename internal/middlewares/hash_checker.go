package middlewares

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

// CheckRequestHashHeader, мидлваря для проверки хэша сообщения на основании заголовка HashSHA256.
func CheckRequestHashHeader(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			hashHeader := r.Header.Get("HashSHA256")
			if key != "" && hashHeader != "" {
				hash := util.Hash(body, key)
				if hashHeader != hash {
					w.WriteHeader(http.StatusBadRequest)
					logger.Log.Infof("key %s hashHeader: %s hash: %s", key, hashHeader, hash)

					return
				}
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

type hashWriter struct {
	http.ResponseWriter
	w   io.Writer
	key string
}

func (hw hashWriter) Write(b []byte) (int, error) {
	hw.Header().Add("HashSHA256", util.Hash(b, hw.key))

	bytesCount, err := hw.w.Write(b)
	return bytesCount, fmt.Errorf("%w", err)
}
func (hw *hashWriter) Close() error {
	if c, ok := hw.w.(io.WriteCloser); ok {
		err := c.Close()
		return fmt.Errorf("%w", err)
	}

	err := errors.New("middlewares: io.WriteCloser is unavailable on the writer")
	return fmt.Errorf("%w", err)
}

// ResposeHeaderWithHash, мидлваря, которая добавляет хэш от контента в заголовок HashSHA256
func ResposeHeaderWithHash(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wr := w

			if key != "" {
				hw := hashWriter{
					ResponseWriter: w,
					w:              w,
					key:            key,
				}

				wr = hw

				defer func() {
					err := hw.Close()
					if err != nil {
						logger.Log.Error("Ошибка закрытия HashWriter, %s", err.Error())
					}
				}()
			}

			next.ServeHTTP(wr, r)
		})
	}
}
