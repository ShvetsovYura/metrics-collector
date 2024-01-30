package logger

import (
	"log/slog"

	"github.com/go-chi/httplog/v2"
)

var HTTPLogger *httplog.Logger

func NewHTTPLogger() {
	HTTPLogger = httplog.NewLogger("metrics-http-logger", httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelInfo,
		Concise:          false,
		RequestHeaders:   false,
		MessageFieldName: "message",
	})
}
