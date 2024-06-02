package agent

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
)

// AgentOptoins хранит информацию настроек запуска
// и выполнения агента сбора метрик
type AgentOptions struct {
	// EndpointAddr: адрес отправки метрик
	EndpointAddr string `env:"ADDRESS"`
	// ReportInterval: интервал отправки метрик на сервер
	ReportInterval int `env:"REPORT_INTERVAL"`
	// PoolInterval: интервал сбора метрик
	PoolInterval int `env:"POLL_INTERVAL"`
	// Key: приватный ключ доступа
	Key string `env:"KEY"`
	// RateLimit: сколько одновременно можно выполнять отправку метрик на сервер
	RateLimit int `env:"RATE_LIMIT"`
}

// ParseArgs  парсит входные аргументы в структуру AgentOptions
// если не переданы - берутся значения по-умолчнаию
func (o *AgentOptions) ParseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "localhost:8080", "server endpoint address")
	flag.IntVar(&o.PoolInterval, "p", 2, "metrics gather interval")
	flag.IntVar(&o.ReportInterval, "r", 10, "interval send metrics to server")
	flag.StringVar(&o.Key, "k", "", "Secret key")
	flag.IntVar(&o.RateLimit, "l", 0, "limit concurent")
	flag.Parse()
}

// ParseEnvs парсит переменные окружения в структуру AgentOptions
func (o *AgentOptions) ParseEnvs() error {
	if err := env.Parse(o); err != nil {
		return errors.New("failed to parse agent env")
	}
	return nil
}
