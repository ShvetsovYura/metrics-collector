package middlewares

import (
	"bytes"
	"io"

	"net/http"

	"github.com/ShvetsovYura/metrics-collector/internal/util"
)

// DecryptMessage: предназначена для расшифровки входящего тела запроса
func DecryptMessage(privateKeyPath string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		handler := func(w http.ResponseWriter, req *http.Request) {
			if privateKeyPath != "" {
				data, errRead := io.ReadAll(req.Body)
				if errRead != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return

				}
				decrytedMessage, err := util.DecryptData(data, privateKeyPath)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				req.Body = io.NopCloser(bytes.NewReader(decrytedMessage))

			}
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(handler)
	}
}
