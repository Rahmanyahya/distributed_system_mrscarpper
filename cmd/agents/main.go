package main

import (
	"bytes"
	"context"
	"distributed_system/internal/config"
	domainConfig "distributed_system/internal/domain/config"
	"distributed_system/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	version int = 0
	countFetch int = 0
	RWMutex sync.RWMutex
)

func main() {
	// Get config path from env or use default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	agentsCfg, err := config.LoadConfigAgents(configPath)
	if err != nil {
		log.Fatalf("Failed to load agents config: %v", err)
	}

	// Self registration
	credential, err := selfRegistration(agentsCfg)
	if err != nil {
		log.Fatalf("[Agent] Failed to read internal key: %v", err)
	}

	log.Println("============================================================")
	log.Println("[Agent] Starting...")
	log.Printf("[Agent] Controller URL: %s", agentsCfg.Controller.URL)
	log.Printf("[Agent] Worker URL: %s", agentsCfg.Worker.URL)
	log.Println("============================================================")

	log.Println("[Agent] Fetching initial config from Controller...")
	initialConfig, err := fetchConfigFromController(agentsCfg, credential)
	if err != nil {
		log.Fatalf("[Agent] Failed to fetch initial config: %v", err)
	}

	log.Printf("[Agent] Initial config: Version=%d, URL=%s", initialConfig.Version, initialConfig.ConfigURL)

	log.Println("[Agent] Pushing initial config to Worker...")
	if err := pushConfigToWorker(agentsCfg, initialConfig); err != nil {
		log.Printf("[Agent] Warning: Failed to push to Worker: %v", err)
	} else {
		log.Println("[Agent] Successfully pushed initial config to Worker!")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poolingInterval := time.Duration(initialConfig.PoolingInterval) * time.Second

	go startPolling(ctx, agentsCfg, initialConfig, credential, poolingInterval)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh
	log.Println("[Agent] Shutting down...")
	cancel()

	log.Println("[Agent] Stopped.")
}

func startPolling(ctx context.Context, agentsCfg *config.ConfigAgents, lastConfig *domainConfig.Config, credential string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("[Agent] Started polling every %v", interval)

	for {
		select {
		case <-ticker.C:
			newConfig, err := fetchConfigFromController(agentsCfg, credential)
			if err != nil {
				log.Printf("[Agent] Error fetching config: %v", err)
				continue
			}

			RWMutex.Lock()
			countFetch += 1

			if newConfig != lastConfig || countFetch > 3{
				countFetch = 0
				log.Printf("[Agent] Config changed! Version: %d -> %d", lastConfig.Version, newConfig.Version)

				if err := pushConfigToWorker(agentsCfg, newConfig); err != nil {
					log.Printf("[Agent] Error pushing to Worker: %v", err)
				} else {
					log.Printf("[Agent] Successfully pushed updated config (version %d) to Worker!", newConfig.Version)
				}

				lastConfig = newConfig
				newInterval := time.Duration(newConfig.PoolingInterval) * time.Second
				if newInterval != interval {
					log.Printf("[Agent] Pooling interval changed: %v -> %v", interval, newInterval)
					ticker.Reset(newInterval)
					interval = newInterval
				}
			} else {
				log.Printf("[Agent] Config unchanged (version %d)", newConfig.Version)
			}
			RWMutex.Unlock()
		case <-ctx.Done():
			log.Println("[Agent] Polling stopped")
			return
		}
	}
}

func fetchConfigFromController(cfg *config.ConfigAgents, credential string) (*domainConfig.Config, error) {
	url := fmt.Sprintf("%s/config/agent", cfg.Controller.URL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+credential)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()


	for resp.StatusCode != http.StatusOK {
		fmt.Println("[Agent] Got non-200 status code from Controller, trying again...")
		time.Sleep(30 * time.Second)
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %w", err)
		}
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var response struct {
		Status string    `json:"status"`
		Data domainConfig.Config `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	RWMutex.Lock()
	if version == 0 || response.Data.Version > version {
		utils.WriteJson("config", &domainConfig.Config{
			Version: response.Data.Version,
			ConfigURL: response.Data.ConfigURL,
			PoolingInterval: response.Data.PoolingInterval,
			UUID: response.Data.UUID,
			CreatedAt: response.Data.CreatedAt,
		})	
	}
	RWMutex.Unlock()

	return &response.Data, nil
}

func selfRegistration(cfg *config.ConfigAgents) (string, error) {
	type Credential struct {
		CredentialKey string `json:"credential_key"`
	}

	// check if already registered
	credential, _ := utils.ReadJSON[Credential]("credential")

	if credential != nil && credential.CredentialKey != "" {
		return credential.CredentialKey, nil
	}

	url := fmt.Sprintf("%s/agent/register", cfg.Controller.URL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "",fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Identity.InternalKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	utils.WriteJson("credential", Credential{
		CredentialKey: response.Data,
	})

	return response.Data, nil
}

func pushConfigToWorker(cfg *config.ConfigAgents, config *domainConfig.Config) error {
	workerConfig := map[string]interface{}{
		"config_url":       config.ConfigURL,
		"pooling_interval": config.PoolingInterval,
		"version":          config.Version,
		"uuid":             config.UUID,
	}

	jsonData, err := json.Marshal(workerConfig)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	url := fmt.Sprintf("%s/config", cfg.Worker.URL)
	fmt.Println(url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.Worker.InternalKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}