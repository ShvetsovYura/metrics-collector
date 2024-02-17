package middlewares

import (
	"net/http"
)

func HashCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashHeader := r.Header.Get("HashSHA256")

		if hashHeader == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
