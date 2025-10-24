package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string     `yaml:"env"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Postgres struct {
	DSN            string        `yaml:"dsn"`
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
		log.Println("CONFIG_PATH not set, using default:", configPath)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	// Read config from YAML
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// Log config
	log.Printf("loaded config from %s\n", configPath)

	return &cfg, nil
}
