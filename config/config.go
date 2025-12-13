package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary          Primary        `koanf:"primary" validate:"required"`
	Server           ServerConfig   `koanf:"server" validate:"required"`
	Database         DatabaseConfig `koanf:"database" validate:"required"`
	Ai               AI             `koanf:"ai" validate:"required"`
	Key              Key            `koanf:"key" validate:"required"`
	Logging          LoggingConfig  `koanf:"logging" validate:"required"`
}

type Primary struct {
	Env          string `koanf:"env" validate:"required,oneof=local prod"`
	ServiceName  string `koanf:"service_name" validate:"required"`
	DatabaseType string `koanf:"database_type" validate:"required,oneof=mock postgres"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

type AI struct {
	LLMInterfaceProvider string `koanf:"llm_provider" validate:"required"`
	LLMInterfaceApiKey   string `koanf:"llm_api_key" validate:"required"`
	VLMInterfaceProvider string `koanf:"vlm_provider" validate:"required"`
	VLMInterfaceApiKey   string `koanf:"vlm_api_key" validate:"required"`
}

type Key struct {
	SecretKey string `koanf:"secret_key" validate:"required"`
}

type LoggingConfig struct {
	Level  string `koanf:"level" validate:"required,oneof=debug info warn error"`
	Format string `koanf:"format" validate:"required,oneof=json text"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	k := koanf.New(".")
	err := k.Load(env.Provider("", k.Delim(), func(s string) string {
		return strings.ToLower(s)
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	config := &Config{}
	if err := k.Unmarshal("", &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config %w", err)
	}

	if err := validator.New().Struct(config); err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal main config")
	}

	return config, nil
}

func (cfg *Config) IsProd() bool {
	return cfg.Primary.Env == "prod"
}
