package agent

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ShvetsovYura/metrics-collector/internal/logger"
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

const (
	EndpointAddrDef   = "localhost:8080"
	ReportIntervalDef = time.Duration(2 * time.Second)
	PollIntervalDef   = time.Duration(2 * time.Second)
)

// AgentOptoins хранит информацию настроек запуска
// и выполнения агента сбора метрик
type Options struct {
	EndpointAddr   string        `env:"ADDRESS" json:"address"`                 // EndpointAddr: адрес отправки метрик
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"` // ReportInterval: интервал отправки метрик на сервер
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`     // PoolInterval: интервал сбора метрик
	Key            string        `env:"KEY"`                                    // Key: приватный ключ доступа
	RateLimit      int           `env:"RATE_LIMIT"`                             // RateLimit: сколько одновременно можно выполнять отправку метрик на сервер
	CryptoKey      string        `env:"CRYPTO_KEY" json:"crypto_key"`           // CryptoKey: путь до файла с публичным ключом
}

func ReadOptions() *Options {
	opt := &Options{}
	opt.parseArgs()
	opt.parseEnvs()
	opt.parseConfig()
	opt.applyDefaultParams()
	return opt
}
func (o *Options) parseConfig() {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "", "path to config file")
	flag.Parse()

	if configPath != "" {
		o.applyConfig(configPath)
		return
	}

	val, ok := os.LookupEnv("CONFIG")
	if ok {
		o.applyConfig(val)
		return
	}
}

func (o *Options) applyDefaultParams() {
	if o.EndpointAddr == "" {
		o.EndpointAddr = EndpointAddrDef
	}
	if o.PollInterval == 0 {
		o.PollInterval = PollIntervalDef
	}
	if o.ReportInterval == 0 {
		o.ReportInterval = ReportIntervalDef
	}

}

// ParseArgs  парсит входные аргументы в структуру AgentOptions
// если не переданы - берутся значения по-умолчнаию
func (o *Options) parseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "", "server endpoint address")
	flag.DurationVar(&o.PollInterval, "p", 0, "metrics gather interval")
	flag.DurationVar(&o.ReportInterval, "r", 0, "interval send metrics to server")
	flag.StringVar(&o.Key, "k", "", "hash key")
	flag.IntVar(&o.RateLimit, "l", 0, "limit concurent")
	flag.StringVar(&o.CryptoKey, "crypto-key", "", "path to public key")

	flag.Parse()
	logger.Log.Infof("flags: %v", *o)
}

// ParseEnvs парсит переменные окружения в структуру AgentOptions
func (o *Options) parseEnvs() error {

	opt := &Options{}
	if err := env.Parse(opt); err != nil {
		return errors.New("failed to parse agent env")
	}

	setParams(o, opt)

	return nil
}

func (o *Options) applyConfig(path string) {
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			logger.Log.Fatal("Не могу открыть конфигурационный файл")
		}
		data, err := io.ReadAll(f)
		if err != nil {
			logger.Log.Fatal("Не могу прочитать конфигурационный файл")
		}

		opt := &Options{}
		json.Unmarshal(data, opt)
		// решение в лоб
		setParams(o, opt)
	}
}
func (o *Options) UnmarshalJSON(data []byte) error {
	type OptionsAlias Options

	optionsValue := &struct {
		*OptionsAlias
		ReportInterval string `json:"report_interval"`
		PollInterval   string `json:"poll_interval"`
	}{
		OptionsAlias: (*OptionsAlias)(o),
	}
	if err := json.Unmarshal(data, optionsValue); err != nil {
		return fmt.Errorf("ошибка парсинга конфигурации %w", err)
	}
	var err error
	o.ReportInterval, err = time.ParseDuration(optionsValue.ReportInterval)
	if err != nil {
		return fmt.Errorf("ошибка преобразования поля ReportInterval %w", err)
	}
	o.PollInterval, err = time.ParseDuration(optionsValue.PollInterval)
	if err != nil {
		return fmt.Errorf("ошибка преобразования поля PollInterval %w", err)
	}
	return nil
}

func setParams(curOpt *Options, tempOpt *Options) {
	if curOpt.EndpointAddr == "" && tempOpt.EndpointAddr != "" {
		curOpt.EndpointAddr = tempOpt.EndpointAddr
	}
	if curOpt.ReportInterval == 0 && tempOpt.ReportInterval != 0 {
		curOpt.ReportInterval = tempOpt.ReportInterval
	}
	if curOpt.PollInterval == 0 && tempOpt.PollInterval != 0 {
		curOpt.PollInterval = tempOpt.PollInterval
	}
	if curOpt.Key == "" && tempOpt.Key != "" {
		curOpt.Key = tempOpt.Key
	}
	if curOpt.RateLimit == 0 && tempOpt.RateLimit != 0 {
		curOpt.RateLimit = tempOpt.RateLimit
	}
	if curOpt.CryptoKey == "" && tempOpt.CryptoKey != "" {
		curOpt.CryptoKey = tempOpt.CryptoKey
	}
}
