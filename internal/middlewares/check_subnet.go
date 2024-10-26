package middlewares

import (
	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/ShvetsovYura/metrics-collector/internal/validator"
)

func CheckTrustetSubnet(trustedSubnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(resp http.ResponseWriter, req *http.Request) {
			xRealIP := req.Header.Get("x-real-ip")
			logger.Log.Debug("входящий ip %s", xRealIP)
			if trustedSubnet != "" {
				if xRealIP == "" {
					resp.WriteHeader(http.StatusForbidden)
					return
				}
				val, err := validator.IsIPInSubnet(xRealIP, trustedSubnet)
				if err != nil {
					resp.WriteHeader(http.StatusInternalServerError)
					resp.Write([]byte(err.Error()))
					return
				}
				if !val {
					resp.WriteHeader(http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(resp, req)
		}
		return http.HandlerFunc(handler)
	}
}
