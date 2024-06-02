package logger

import (
	"log/slog"

	"github.com/go-chi/httplog/v2"
)

// HTTPLogger, глобальный http-логгер приложения
var HTTPLogger *httplog.Logger

// NewHTTPLogger, создание логгера http-запросов.
func NewHTTPLogger() {
	HTTPLogger = httplog.NewLogger("metrics-http-logger", httplog.Options{
		JSON:             true,
		LogLevel:         slog.LevelInfo,
		Concise:          false,
		RequestHeaders:   false,
		MessageFieldName: "message",
	})
}
