package middlewares

import (
	"net/http"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/custom"
	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

func WithLog(h http.Handler) http.HandlerFunc {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		data := &custom.RespData{
			StatusCode: 0,
			Size:       0,
		}

		lw := custom.LogResponseWriter{
			ResponseWriter: w,
			Data:           data,
		}

		h.ServeHTTP(&lw, r)
		duration := time.Since(start)

		logger.Log.Infoln(
			"status", data.StatusCode,
			"size", data.Size,
			"method", r.Method,
			"url", r.RequestURI,
			"duration:", duration,
		)
	}
	return http.HandlerFunc(logFunc)
}
