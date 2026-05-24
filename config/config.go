package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config contains application configuration.
type Config struct {
	Env        string     `yaml:"env"`
	HTTPServer HTTPServer `yaml:"http_server"`
	GRPCServer GRPCServer `yaml:"grpc_server"`
	Postgres   Postgres   `yaml:"postgres"`
	Redis      Redis      `yaml:"redis"`
}

// HTTPServer contains HTTP server configuration.
type HTTPServer struct {
	Port         string        `yaml:"port" env-default:":8080"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// GRPCServer contains gRPC server configuration.
type GRPCServer struct {
	Port string `yaml:"port" env-default:":9090"`
}

// Postgres contains PostgreSQL configuration.
type Postgres struct {
	DSN            string        `yaml:"dsn"`
	MaxRetries     int           `yaml:"max_retries" env-default:"10"`
	RetryInterval  time.Duration `yaml:"retry_interval" env-default:"5s"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" env-default:"30s"`
}

// Redis contains Redis configuration.
type Redis struct {
	Addr           string        `yaml:"addr" env-default:"localhost:6379"`
	Password       string        `yaml:"password"`
	DB             int           `yaml:"db"`
	TTL            time.Duration `yaml:"ttl" env-default:"5m"`
	MaxRetries     int           `yaml:"max_retries" env-default:"10"`
	RetryInterval  time.Duration `yaml:"retry_interval" env-default:"5s"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" env-default:"30s"`
}

// New returns a new config.
func New() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("config file not available %q: %w", configPath, err)
	}

	// Read config from YAML
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("read config %q: %w", configPath, err)
	}

	return &cfg, nil
}
