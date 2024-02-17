package agent

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
)

type AgentOptions struct {
	EndpointAddr   string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
}

func (o *AgentOptions) ParseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "server endpoint address")
	flag.IntVar(&o.PoolInterval, "p", 2, "metrics gather interval")
	flag.IntVar(&o.ReportInterval, "r", 10, "interval send metrics to server")
	flag.StringVar(&o.Key, "k", "", "Secret key")
	flag.Parse()
}
func (o *AgentOptions) ParseEnvs() error {
	if err := env.Parse(o); err != nil {
		return errors.New("failed to parse agent env")
	}
	return nil
}
