package custom

import "net/http"

type (
	RespData struct {
		StatusCode int
		Size       int
	}

	LogResponseWriter struct {
		http.ResponseWriter
		Data *RespData
	}
)

func (w *LogResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.Data.Size += size
	return size, err
}

func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.Data.StatusCode = statusCode
}
