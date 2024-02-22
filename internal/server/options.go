package server

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
)

type ServerOptions struct {
	EndpointAddr    string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DBDSN           string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
}

func (o *ServerOptions) ParseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "endpoint address")
	flag.IntVar(&o.StoreInterval, "i", 300, "interval to store data on file. 0 for immediately")
	flag.StringVar(&o.FileStoragePath, "f", "/tmp/metrics-db.json", "path to save metrics values")
	flag.BoolVar(&o.Restore, "r", true, "restoring metrics values on start")
	flag.StringVar(&o.DBDSN, "d", "", "database connection DSN")
	flag.StringVar(&o.Key, "k", "", "Secret key value")

	flag.Parse()
}

func (o *ServerOptions) ParseEnvs() error {
	if err := env.Parse(o); err != nil {
		return errors.New("failed to parse server env")
	}
	return nil
}
