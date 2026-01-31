package scheduler

import (
	"context"
	"distributed_system/internal/agent/client"
	"distributed_system/internal/domain/config"
	"fmt"
	"log"
	"sync"
	"time"
)

// ConfigScheduler handles periodic config fetching from the controller
type ConfigScheduler struct {
	controllerURL string
	token         string
	client        *client.ConfigClient
	currentConfig *config.Config
	mu            sync.RWMutex
	ticker        *time.Ticker
	tickerMu      sync.Mutex
	stopCh        chan struct{}
	running       bool
	onConfigUpdate func(*config.Config)
}

func NewConfigScheduler(controllerURL, token string, onConfigUpdate func(*config.Config)) *ConfigScheduler {
	return &ConfigScheduler{
		controllerURL:  controllerURL,
		token:          token,
		client:         client.NewConfigClient(controllerURL),
		stopCh:         make(chan struct{}),
		onConfigUpdate: onConfigUpdate,
	}
}

// Start begins the periodic config fetching with initial interval
func (s *ConfigScheduler) Start(ctx context.Context, initialInterval int) {
	s.tickerMu.Lock()
	if s.running {
		s.tickerMu.Unlock()
		log.Println("[ConfigScheduler] Already running")
		return
	}
	s.running = true
	s.tickerMu.Unlock()

	interval := time.Duration(initialInterval) * time.Second
	s.tickerMu.Lock()
	s.ticker = time.NewTicker(interval)
	s.tickerMu.Unlock()

	log.Printf("[ConfigScheduler] Started. Checking config every %d seconds from %s",
		initialInterval, s.controllerURL)

	// Initial fetch
	s.fetchConfig(ctx)

	// Periodic fetch
	go func() {
		for {
			select {
			case <-s.getTicker():
				s.fetchConfig(ctx)
			case <-s.stopCh:
				log.Println("[ConfigScheduler] Stopped")
				return
			case <-ctx.Done():
				log.Println("[ConfigScheduler] Context cancelled")
				return
			}
		}
	}()
}

// getTicker safely returns the ticker channel
func (s *ConfigScheduler) getTicker() <-chan time.Time {
	s.tickerMu.Lock()
	defer s.tickerMu.Unlock()
	if s.ticker != nil {
		return s.ticker.C
	}
	return nil
}

// Stop stops the scheduler
func (s *ConfigScheduler) Stop() {
	s.tickerMu.Lock()
	defer s.tickerMu.Unlock()

	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
	s.running = false

	select {
	case <-s.stopCh:
		// Already closed
	default:
		close(s.stopCh)
	}
}

// GetConfig returns the current config (thread-safe)
func (s *ConfigScheduler) GetConfig() *config.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentConfig
}

// updateInterval updates the ticker with new interval
func (s *ConfigScheduler) updateInterval(newInterval int) {
	s.tickerMu.Lock()
	defer s.tickerMu.Unlock()

	if s.ticker != nil {
		s.ticker.Stop()
	}

	interval := time.Duration(newInterval) * time.Second
	s.ticker = time.NewTicker(interval)

	log.Printf("[ConfigScheduler] Interval updated to %d seconds", newInterval)
}

// fetchConfig fetches the latest config from the controller
func (s *ConfigScheduler) fetchConfig(ctx context.Context) {
	log.Printf("[ConfigScheduler] Fetching config from controller...")

	newConfig, err := s.client.GetLatestConfig(ctx, s.token)
	if err != nil {
		log.Printf("[ConfigScheduler] Error fetching config: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if config has changed
	if s.currentConfig == nil || newConfig.Version > s.currentConfig.Version {
		log.Printf("[ConfigScheduler] New config received! Version: %d (was %d)",
			newConfig.Version, func() int {
				if s.currentConfig != nil {
					return s.currentConfig.Version
				}
				return 0
			}())

		// Check if pooling interval changed
		if s.currentConfig != nil && s.currentConfig.PoolingInterval != newConfig.PoolingInterval {
			log.Printf("[ConfigScheduler] Pooling interval changed: %d -> %d",
				s.currentConfig.PoolingInterval, newConfig.PoolingInterval)
			s.updateInterval(newConfig.PoolingInterval)
		}

		s.currentConfig = newConfig

		// Call callback if config updated
		if s.onConfigUpdate != nil {
			s.onConfigUpdate(newConfig)
		}
	} else {
		log.Printf("[ConfigScheduler] Config unchanged. Version: %d", newConfig.Version)
	}
}

// ForceFetch forces an immediate config fetch
func (s *ConfigScheduler) ForceFetch(ctx context.Context) error {
	log.Println("[ConfigScheduler] Force fetching config...")

	newConfig, err := s.client.GetLatestConfig(ctx, s.token)
	if err != nil {
		return fmt.Errorf("failed to force fetch config: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentConfig = newConfig

	// Call callback if config updated
	if s.onConfigUpdate != nil {
		s.onConfigUpdate(newConfig)
	}

	return nil
}

// SetInitialConfig sets the initial config and starts the scheduler
func (s *ConfigScheduler) SetInitialConfig(ctx context.Context, cfg *config.Config) {
	s.mu.Lock()
	s.currentConfig = cfg
	s.mu.Unlock()

	log.Printf("[ConfigScheduler] Initial config set. Version: %d, Pooling Interval: %d seconds",
		cfg.Version, cfg.PoolingInterval)

	// Start scheduler with initial interval
	s.Start(ctx, cfg.PoolingInterval)
}
