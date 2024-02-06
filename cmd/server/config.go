package main

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
)

type ServerOptions struct {
	EndpointAddr string `env:"ADDRESS"`
	DBDSN        string `env:"DATABASE_DSN"`
}

func (o *ServerOptions) ParseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "endpoint address")
	flag.StringVar(&o.DBDSN, "d", "", "database connection DSN")
	flag.Parse()
}

func (o *ServerOptions) ParseEnvs() error {
	if err := env.Parse(o); err != nil {
		return errors.New("failed to parse server env")
	}
	return nil
}
