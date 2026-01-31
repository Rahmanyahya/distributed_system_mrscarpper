package version_checker

import (
	"bytes"
	"context"
	"distributed_system/internal/domain/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// VersionConfig represents the version response from Controller
type VersionConfig struct {
	Version int `json:"version"`
}

// VersionChecker checks for config version changes from Controller's Redis
type VersionChecker struct {
	controllerURL    string
	controllerToken  string
	workerURL        string
	client           *http.Client
	currentVersion   int
	mu               sync.RWMutex
	ticker           *time.Ticker
	tickerMu         sync.Mutex
	stopCh           chan struct{}
	running          bool
	onConfigUpdate   func(*config.Config)
}

// NewVersionChecker creates a new version checker
func NewVersionChecker(controllerURL, controllerToken, workerURL string, onConfigUpdate func(*config.Config)) *VersionChecker {
	return &VersionChecker{
		controllerURL:    controllerURL,
		controllerToken:  controllerToken,
		workerURL:        workerURL,
		client:           &http.Client{Timeout: 10 * time.Second},
		currentVersion:   0,
		stopCh:           make(chan struct{}),
		onConfigUpdate:   onConfigUpdate,
	}
}

// Start begins the periodic version checking
func (vc *VersionChecker) Start(ctx context.Context, checkInterval int) {
	vc.tickerMu.Lock()
	if vc.running {
		vc.tickerMu.Unlock()
		log.Println("[VersionChecker] Already running")
		return
	}
	vc.running = true
	vc.tickerMu.Unlock()

	interval := time.Duration(checkInterval) * time.Second
	vc.tickerMu.Lock()
	vc.ticker = time.NewTicker(interval)
	vc.tickerMu.Unlock()

	log.Printf("[VersionChecker] Started. Checking version every %d seconds", checkInterval)

	// Initial check
	vc.checkVersion(ctx)

	// Periodic check
	go func() {
		for {
			select {
			case <-vc.getTicker():
				vc.checkVersion(ctx)
			case <-vc.stopCh:
				log.Println("[VersionChecker] Stopped")
				return
			case <-ctx.Done():
				log.Println("[VersionChecker] Context cancelled")
				return
			}
		}
	}()
}

// Stop stops the version checker
func (vc *VersionChecker) Stop() {
	vc.tickerMu.Lock()
	defer vc.tickerMu.Unlock()

	if vc.ticker != nil {
		vc.ticker.Stop()
		vc.ticker = nil
	}
	vc.running = false

	select {
	case <-vc.stopCh:
		// Already closed
	default:
		close(vc.stopCh)
	}
}

// SetInitialVersion sets the initial version
func (vc *VersionChecker) SetInitialVersion(version int) {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vc.currentVersion = version
	log.Printf("[VersionChecker] Initial version set to %d", version)
}

// getTicker safely returns the ticker channel
func (vc *VersionChecker) getTicker() <-chan time.Time {
	vc.tickerMu.Lock()
	defer vc.tickerMu.Unlock()
	if vc.ticker != nil {
		return vc.ticker.C
	}
	return nil
}

// checkVersion checks if the version has changed in Controller's Redis
func (vc *VersionChecker) checkVersion(ctx context.Context) {
	log.Printf("[VersionChecker] Checking version from Controller Redis...")

	// Check version from Controller
	remoteVersion, err := vc.fetchVersionFromController(ctx)
	if err != nil {
		log.Printf("[VersionChecker] Error fetching version: %v", err)
		return
	}

	vc.mu.RLock()
	localVersion := vc.currentVersion
	vc.mu.RUnlock()

	// Check if version changed
	if remoteVersion > localVersion {
		log.Printf("[VersionChecker] Version changed! Local: %d, Remote: %d", localVersion, remoteVersion)

		// Fetch full config from Controller
		newConfig, err := vc.fetchConfigFromController(ctx)
		if err != nil {
			log.Printf("[VersionChecker] Error fetching config: %v", err)
			return
		}

		// Update local version
		vc.mu.Lock()
		vc.currentVersion = remoteVersion
		vc.mu.Unlock()

		// Push config to Worker
		if err := vc.pushConfigToWorker(newConfig); err != nil {
			log.Printf("[VersionChecker] Error pushing config to Worker: %v", err)
			return
		}

		// Call callback if set
		if vc.onConfigUpdate != nil {
			vc.onConfigUpdate(newConfig)
		}

		log.Printf("[VersionChecker] Successfully updated Worker to version %d", remoteVersion)
	} else {
		log.Printf("[VersionChecker] Version unchanged: %d", localVersion)
	}
}

// fetchVersionFromController fetches the current version from Controller's Redis
func (vc *VersionChecker) fetchVersionFromController(ctx context.Context) (int, error) {
	url := fmt.Sprintf("%s/config/version", vc.controllerURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+vc.controllerToken)

	resp, err := vc.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Data VersionConfig `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Data.Version, nil
}

// fetchConfigFromController fetches the full config from Controller
func (vc *VersionChecker) fetchConfigFromController(ctx context.Context) (*config.Config, error) {
	url := fmt.Sprintf("%s/config/agent", vc.controllerURL)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+vc.controllerToken)

	resp, err := vc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Data config.Config `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response.Data, nil
}

// pushConfigToWorker pushes the new configuration to Worker
func (vc *VersionChecker) pushConfigToWorker(cfg *config.Config) error {
	workerConfig := map[string]interface{}{
		"config_url":       cfg.ConfigURL,
		"pooling_interval": cfg.PoolingInterval,
		"version":          cfg.Version,
		"uuid":             cfg.UUID,
	}

	jsonData, err := json.Marshal(workerConfig)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	url := fmt.Sprintf("%s/config", vc.workerURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := vc.client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Printf("[VersionChecker] Successfully pushed config to Worker at %s", vc.workerURL)
	return nil
}
