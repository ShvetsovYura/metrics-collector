package util

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AgentOptions struct {
	EndpointAddr   string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
}

type ServerOptions struct {
	EndpointAddr string `env:"ADDRESS"`
}

func (o *AgentOptions) ParseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "server endpoint address")
	flag.IntVar(&o.PoolInterval, "p", 2, "metrics gather interval")
	flag.IntVar(&o.ReportInterval, "r", 10, "interval send metrics to server")
	flag.Parse()
}
func (o *AgentOptions) ParseEnvs() {
	env.Parse(o)
}

func (o *ServerOptions) ParseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "endpoint address")
	flag.Parse()
}

func (o *ServerOptions) ParseEnvs() {
	if err := env.Parse(o); err != nil {
		fmt.Println("ERROR ", err)
	}

}
