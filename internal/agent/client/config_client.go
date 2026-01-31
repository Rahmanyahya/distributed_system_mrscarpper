package client

import (
	"context"
	"distributed_system/internal/domain/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ConfigClient handles HTTP communication with the controller
type ConfigClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewConfigClient creates a new config client
func NewConfigClient(baseURL string) *ConfigClient {
	return &ConfigClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetLatestConfig fetches the latest config from the controller
func (c *ConfigClient) GetLatestConfig(ctx context.Context, token string) (*config.Config, error) {
	url := fmt.Sprintf("%s/config/agent", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header with bearer token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response - expecting {"data": {...}} format from response.Success
	var response struct {
		Data config.Config `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response.Data, nil
}
