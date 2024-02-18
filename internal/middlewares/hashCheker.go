package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

func CheckHashHeader(key string) func(next http.Handler) http.Handler {
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
					w.WriteHeader(http.StatusTeapot)
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
