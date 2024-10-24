package server

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
	ServerTypeDef    = "grpc"
	EndpointAddrDef  = "localhost:8080"
	StoreIntervalDef = time.Duration(300 * time.Second)
	RestoreDef       = true
	LogLevelDef      = "info"
)

// ServerOptions, хранит опции сервера сбора метрик.
type Options struct {
	ServerType      string        `env:"SERVR_TYPE" json:"server_type"`        // ServerType: тип запускаемого сервера метрик (http, grpc)
	EndpointAddr    string        `env:"ADDRESS" json:"address"`               // адрес запуска сервера сбора метрик
	StoreInterval   time.Duration `env:"STORE_INTERVAL" json:"store_interval"` // интервал сохранения метрик в хранилище
	FileStoragePath string        `env:"STORE_FILE" json:"store_file"`         // путь до сохранения метрик в файл
	Restore         bool          `env:"RESTORE" json:"restore"`               // восстанавливать метрики при старте приложения
	DBDSN           string        `env:"DATABASE_DSN" json:"database_dsn"`     // строка подключения к БД
	Key             string        `env:"KEY" json:"key"`                       // ключ хеширования сообщения
	CryptoKey       string        `env:"CRYPTO_KEY" json:"crypto_key"`         // путь до файла с приватным ключом
	TrustedSubnet   string        `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	LogLevel        string        `env:"LOG_LEVEL" json:"log_level"`
}

func ReadOptions() *Options {
	opt := &Options{}
	// чтение аругментов
	opt.parseArgs()
	// чтение переменных окружения и перезапись пустых значений
	if err := opt.parseEnvs(); err != nil {
		logger.Log.Fatal(err.Error())
	}
	// чтение конфига и перезапись пустых значений
	if path := opt.getConfigPath(); path != "" {
		opt.applyConfig(path)
	}
	// применение дефолтных значений в оставшиеся пустые пераметры
	opt.applyDefaultParams()
	return opt
}

func (o *Options) UnmarshalJSON(data []byte) error {
	type OptionsAlias Options

	optionsValue := &struct {
		*OptionsAlias
		StoreInterval string `json:"store_interval"`
	}{
		OptionsAlias: (*OptionsAlias)(o),
	}
	if err := json.Unmarshal(data, optionsValue); err != nil {
		return fmt.Errorf("ошибка парсинга конфигурации %w", err)
	}
	var err error
	o.StoreInterval, err = time.ParseDuration(optionsValue.StoreInterval)
	if err != nil {
		return fmt.Errorf("ошибка преобразования поля StoreInterval %w", err)
	}

	return nil
}

func (o *Options) getConfigPath() string {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "", "path to config file")
	flag.Parse()

	if configPath != "" {
		return configPath
	}

	val, ok := os.LookupEnv("CONFIG")
	if ok {
		return val
	}
	return ""
}

func (o *Options) applyDefaultParams() {
	if o.ServerType == "" {
		o.ServerType = ServerTypeDef
	}
	if o.EndpointAddr == "" {
		o.EndpointAddr = EndpointAddrDef
	}
	if o.StoreInterval == -1 {
		o.StoreInterval = StoreIntervalDef
	}
	if !o.Restore {
		o.Restore = RestoreDef
	}
	if o.LogLevel == "" {
		o.LogLevel = LogLevelDef
	}
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
		// не очень красиво
		reassignOptions(o, opt)
	}
}

func (o *Options) parseEnvs() error {
	opt := &Options{}
	if err := env.Parse(opt); err != nil {
		return errors.New("failed to parse agent env")
	}
	reassignOptions(o, opt)
	return nil
}

// ParseArgs, парсит значения аргументов в опции сервера сбора метрик.
func (o *Options) parseArgs() {
	flag.StringVar(&o.EndpointAddr, "a", "", "endpoint address")
	flag.DurationVar(&o.StoreInterval, "i", -1, "interval to store data on file. 0 for immediately")
	flag.StringVar(&o.FileStoragePath, "f", "/tmp/metrics-db.json", "path to save metrics values")
	flag.BoolVar(&o.Restore, "r", false, "restoring metrics values on start")
	flag.StringVar(&o.DBDSN, "d", "", "database connection DSN")
	flag.StringVar(&o.Key, "k", "", "Secret key value")
	flag.StringVar(&o.CryptoKey, "crypto-key", "", "path to private key")
	flag.StringVar(&o.TrustedSubnet, "t", "", "verify client in trusted subnet")

	flag.Parse()
}

// ParseEnvs, парсит значения из переменных окружения в опции сервера сбора метрик.
func (o *Options) ParseEnvs() error {
	opt := &Options{}
	if err := env.Parse(opt); err != nil {
		return errors.New("failed to parse server env")
	}
	reassignOptions(o, opt)
	return nil
}

// reassignOptions переопределяет значения опции
func reassignOptions(curOpt *Options, tempOpt *Options) {
	// установка только если значение параметров пустое, а целевое - нет
	if curOpt.EndpointAddr == "" && tempOpt.EndpointAddr != "" {
		curOpt.EndpointAddr = tempOpt.EndpointAddr
	}
	if curOpt.StoreInterval == -1 && tempOpt.StoreInterval != 0 {
		curOpt.StoreInterval = tempOpt.StoreInterval
	}
	if curOpt.FileStoragePath == "" && tempOpt.FileStoragePath != "" {
		curOpt.FileStoragePath = tempOpt.FileStoragePath
	}
	// сомнительное присвоение
	if !curOpt.Restore && tempOpt.Restore {
		curOpt.Restore = tempOpt.Restore
	}
	if curOpt.Key == "" && tempOpt.Key != "" {
		curOpt.Key = tempOpt.Key
	}
	if curOpt.DBDSN == "" && tempOpt.DBDSN != "" {
		curOpt.DBDSN = tempOpt.DBDSN
	}
	if curOpt.CryptoKey == "" && tempOpt.CryptoKey != "" {
		curOpt.CryptoKey = tempOpt.CryptoKey
	}
	if curOpt.LogLevel == "" && tempOpt.LogLevel != "" {
		curOpt.LogLevel = tempOpt.LogLevel
	}
	if curOpt.TrustedSubnet == "" && tempOpt.TrustedSubnet != "" {
		curOpt.TrustedSubnet = tempOpt.TrustedSubnet
	}
}
