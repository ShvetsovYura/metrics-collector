package util

import (
	"flag"
)

type AgentOptions struct {
	endpointAddr   string
	reportInterval int
	poolInterval   int
}

func (o *AgentOptions) GetEndpoint() string {
	return o.endpointAddr
}
func (o *AgentOptions) GetReportInterval() int {
	return o.reportInterval
}

func (o *AgentOptions) GetPoolInterval() int {
	return o.poolInterval
}

type ServerOptions struct {
	endpointAddr string
}

func (o *AgentOptions) ParseArgs() {
	flag.StringVar(&o.endpointAddr, "a", "localhost:8080", "server endpoint address")
	flag.IntVar(&o.poolInterval, "p", 2, "metrics gather interval")
	flag.IntVar(&o.reportInterval, "r", 10, "interval send metrics to server")
	flag.Parse()
}

func (o *ServerOptions) ParseArgs() {
	flag.StringVar(&o.endpointAddr, "a", "localhost:8080", "endpoint address")
	flag.Parse()
}

func (o *ServerOptions) GetEndpoint() string {
	return o.endpointAddr
}
