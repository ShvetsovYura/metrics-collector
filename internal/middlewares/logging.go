package middlewares

import (
	"net/http"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
)

func WithLog(h http.Handler) http.HandlerFunc {
	logFunc := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(start)

		logger.Log.Infoln(
			"method", r.Method,
			"url", r.RequestURI,
			"duration:", duration,
		)
	}
	return http.HandlerFunc(logFunc)
}
