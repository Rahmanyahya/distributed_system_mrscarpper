package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type IdentityConfig struct {
	InternalKey string `mapstructure:"internal_key"`
}

type Controller struct {
	URL string `mapstructure:"url"`
}

type Worker struct {
	URL         string `mapstructure:"url"`
	InternalKey string `mapstructure:"internal_key"`
}

type ConfigAgents struct {
	Identity   IdentityConfig `mapstructure:"identity"`
	Controller Controller     `mapstructure:"controller"`
	Worker     Worker         `mapstructure:"worker"`
}

func LoadConfigAgents(path string) (*ConfigAgents, error) {
	v := viper.New()

	v.SetConfigName("agent-config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg ConfigAgents
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}