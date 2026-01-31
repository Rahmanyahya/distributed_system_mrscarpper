package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type WorkerConfig struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Auth struct {
		InternalKey string `mapstructure:"internal_key"`
	} `mapstructure:"auth"`
}

func LoadWorkerConfig(configPath string) (*WorkerConfig, error) {
	viper.SetConfigName("worker-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", 8082)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read worker config file: %w", err)
	}

	var cfg WorkerConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal worker config: %w", err)
	}

	fmt.Printf("Worker config loaded successfully:\n")
	fmt.Printf("  Server Port: %d\n", cfg.Server.Port)

	return &cfg, nil
}
