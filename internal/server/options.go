package server

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
)

// ServerOptions, хранит опции сервера сбора метрик.
type ServerOptions struct {
	EndpointAddr    string `env:"ADDRESS"`           // адрес запуска сервера сбора метрик
	StoreInterval   int    `env:"STORE_INTERVAL"`    // интервал сохранения метрик в хранилище
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // путь до сохранения метрик в файл
	Restore         bool   `env:"RESTORE"`           // восстанавливать метрики при старте приложения
	DBDSN           string `env:"DATABASE_DSN"`      // строка подключения к БД
	Key             string `env:"KEY"`               // ключ хеширования сообщения
}

// ParseArgs, парсит значения аргументов в опции сервера сбора метрик.
func (o *ServerOptions) ParseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "endpoint address")
	flag.IntVar(&o.StoreInterval, "i", 300, "interval to store data on file. 0 for immediately")
	flag.StringVar(&o.FileStoragePath, "f", "/tmp/metrics-db.json", "path to save metrics values")
	flag.BoolVar(&o.Restore, "r", true, "restoring metrics values on start")
	flag.StringVar(&o.DBDSN, "d", "", "database connection DSN")
	flag.StringVar(&o.Key, "k", "", "Secret key value")

	flag.Parse()
}

// ParseEnvs, парсит значения из переменных окружения в опции сервера сбора метрик.
func (o *ServerOptions) ParseEnvs() error {
	if err := env.Parse(o); err != nil {
		return errors.New("failed to parse server env")
	}
	return nil
}
