package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Log, глобальный логер пирложения.
var Log *zap.SugaredLogger = zap.NewNop().Sugar()

// InitLogger, инициализатор логгера приложения.
func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("ошибка получения уровня логирования, %w", err)
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("ошибка создания логера, %w", err)
	}

	Log = zl.Sugar()

	return nil
}
