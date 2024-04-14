package config

import (
	"github.com/spf13/viper"
	filepath2 "path/filepath"
	"runtime"
	"time"
)

const (
	defaultHttpPort    = "8080"
	defaultRWTimeout   = 10 * time.Second
	defaultLoggerLevel = 5
)

type Config struct {
	PostgreSQL PostgreSQLConfig
	HTTP       HTTPConfig
	Logger     LoggerConfig
	JWT        JWTConfig
	Cache      RedisConfig
}

type LoggerConfig struct {
	Level      int
	FileName   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PostgreSQLConfig struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	DBName                string
	SSLMode               string
	MaxIdleConnections    int
	MaxOpenConnections    int
	ConnectionMaxLifeTime time.Duration
	DriverName            string
}

type JWTConfig struct {
	SigningKey string
}

type RedisConfig struct {
	Host               string
	Port               string
	DB                 int
	CacheTTL           time.Duration
	RetryInterval      time.Duration
	MaxNumberOfRetries int
}

func Init(path string) (*Config, error) {
	// setDefault()

	var cfg Config

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseConfigFile(filepath string) error {
	if os := runtime.GOOS; os == "linux" {
		viper.SetConfigFile("/app/configs/main.yaml")
	} else {
		path := filepath2.Dir(filepath)
		name := filepath2.Base(filepath)

		viper.AddConfigPath(path)
		viper.SetConfigName(name)
	}

	return viper.ReadInConfig()
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("logger", &cfg.Logger); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("postgresql", &cfg.PostgreSQL); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("jwt", &cfg.JWT); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("redis", &cfg.Cache); err != nil {
		return err
	}

	return nil
}
