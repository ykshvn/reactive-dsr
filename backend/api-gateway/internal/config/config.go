package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Services  ServicesConfig
	RateLimit RateLimitConfig
	Cors      CorsConfig
}

type ServerConfig struct {
	Port                  int           `mapstructure:"port"`
	ReadTimeout           time.Duration `mapstructure:"read_timeout"`
	WriteTimeout          time.Duration `mapstructure:"write_timeout"`
	IdleTimeout           time.Duration `mapstructure:"idle_timeout"`
	MaxIdleConns          int           `mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost   int           `mapstructure:"max_idle_conns_per_host"`
	ResponseHeaderTimeout time.Duration `mapstructure:"response_header_timeout"`
}

type ServicesConfig struct {
	SimulationServiceURL string `mapstructure:"simulation_url"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

type CorsConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read config: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
