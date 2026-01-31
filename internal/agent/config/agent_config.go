package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// AgentConfig represents the configuration for an agent stored in JSON file
type AgentConfig struct {
	AgentID    string `json:"agent_id"`
	Token      string `json:"token"`
	ControllerURL string `json:"controller_url"`
	CheckInterval int `json:"check_interval"` // in seconds
}

// LoadAgentConfig loads agent configuration from a JSON file
func LoadAgentConfig(path string) (*AgentConfig, error) {
	if path == "" {
		path = "agent_config.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read agent config file: %w", err)
	}

	var cfg AgentConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse agent config file: %w", err)
	}

	// Validate config
	if cfg.AgentID == "" {
		return nil, fmt.Errorf("agent_id is required")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if cfg.ControllerURL == "" {
		return nil, fmt.Errorf("controller_url is required")
	}
	if cfg.CheckInterval < 30 {
		cfg.CheckInterval = 30 // minimum 30 seconds
	}

	return &cfg, nil
}

// SaveAgentConfig saves agent configuration to a JSON file
func SaveAgentConfig(path string, cfg *AgentConfig) error {
	if path == "" {
		path = "agent_config.json"
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal agent config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write agent config file: %w", err)
	}

	return nil
}
