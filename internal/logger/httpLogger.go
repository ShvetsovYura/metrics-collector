package logger

import (
	"log/slog"

	"github.com/go-chi/httplog/v2"
)

var HttpLogger *httplog.Logger

func NewHttpLogger() {
	HttpLogger = httplog.NewLogger("metrics-http-logger", httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelInfo,
		Concise:          false,
		RequestHeaders:   false,
		MessageFieldName: "message",
	})
}
