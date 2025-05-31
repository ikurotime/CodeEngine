package config

import (
	"fmt"
	"ikurotime/code-engine/pkg"
	"os"

	"github.com/go-playground/validator/v10"
)

type ServerConfig struct {
	Port                    string `yaml:"port"`
	MaxConcurrentExecutions int    `yaml:"maxConcurrentExecutions"`
	ExecutionTimeout        int    `yaml:"executionTimeout"`
}

type ContainerConfig struct {
	CPULimit    float64 `yaml:"cpuLimit"`
	MemoryLimit int     `yaml:"memoryLimit"`
}

type DatabaseConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Container ContainerConfig `yaml:"container"`
	Database  DatabaseConfig  `yaml:"database"`
}

func LoadConfig() (*Config, error) {
	var err error

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	cfg := &Config{}

	err = pkg.ReadFile(fmt.Sprintf("config/.env.%s.yaml", env), cfg)
	if err != nil {
		return nil, err
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
