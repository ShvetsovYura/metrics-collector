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
	ClientTypeDef     = "grpc"
	EndpointAddrDef   = "localhost:8080"
	ReportIntervalDef = time.Duration(2 * time.Second)
	PollIntervalDef   = time.Duration(2 * time.Second)
	LogLevelDef       = "debug"
)

// AgentOptoins хранит информацию настроек запуска
// и выполнения агента сбора метрик
type Options struct {
	ClientType     string        `env:"CLIENT_TYPE" json:"client_type"`         // ClientType: тип клиента отправки метрик (http, grpc)
	EndpointAddr   string        `env:"ADDRESS" json:"address"`                 // EndpointAddr: адрес отправки метрик
	ReportInterval time.Duration `env:"REPORT_INTERVAL" json:"report_interval"` // ReportInterval: интервал отправки метрик на сервер
	PollInterval   time.Duration `env:"POLL_INTERVAL" json:"poll_interval"`     // PoolInterval: интервал сбора метрик
	Key            string        `env:"KEY"`                                    // Key: ключ хэширования
	RateLimit      int           `env:"RATE_LIMIT"`                             // RateLimit: сколько одновременно можно выполнять отправку метрик на сервер
	CryptoKey      string        `env:"CRYPTO_KEY" json:"crypto_key"`           // CryptoKey: путь до файла с публичным ключом
	LogLevel       string        `json:"log_level"`
}

func ReadOptions() *Options {
	opt := &Options{}
	// чтение аругментов
	opt.parseArgs()
	// чтение переменных окружения и перезапись пустых значений
	if err := opt.parseEnvs(); err != nil {
		logger.Log.Error(err.Error())
	}

	// чтение конфига и перезапись пустых значений
	opt.parseConfig()
	// применение дефолтных значений в оставшиеся пустые пераметры
	opt.applyDefaultParams()
	return opt
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
	if o.ClientType == "" {
		o.ClientType = ClientTypeDef
	}
	if o.EndpointAddr == "" {
		o.EndpointAddr = EndpointAddrDef
	}
	if o.PollInterval == 0 {
		o.PollInterval = PollIntervalDef
	}
	if o.ReportInterval == 0 {
		o.ReportInterval = ReportIntervalDef
	}
	if o.LogLevel == "" {
		o.LogLevel = LogLevelDef
	}
}

// ParseArgs  парсит входные аргументы в структуру AgentOptions
// если не переданы - берутся значения по-умолчнаию
func (o *Options) parseArgs() {
	// устанавливаем дефолтные значения аргументов в null-значения
	// нужно для дальнейшего переопределения значений
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
	reassignOptions(o, opt)
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
		if err := json.Unmarshal(data, opt); err != nil {
			logger.Log.Error(err.Error())
			return
		}
		// решение в лоб
		reassignOptions(o, opt)
	}
}

// reassignOptions переопределяет значения опции
func reassignOptions(curOpt *Options, tempOpt *Options) {
	// установка только если значение параметров пустое, а целевое - нет
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
	if curOpt.LogLevel == "" && tempOpt.LogLevel != "" {
		curOpt.LogLevel = tempOpt.LogLevel
	}
}
