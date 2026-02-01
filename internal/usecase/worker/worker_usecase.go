package worker

import (
	"context"
	"distributed_system/internal/domain/worker"
	"distributed_system/pkg/errors"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	globalConfig *worker.WorkerConfig
	configMutex  sync.RWMutex
)

type Worker struct {
	httpClient *http.Client
}

func NewWorkerUsecase(httpClient *http.Client) worker.Usecase {
	return &Worker{httpClient: httpClient}
}

func (u *Worker) Hit(ctx context.Context) (any, error) {
	configMutex.Lock()
	if globalConfig == nil {
		return nil, errors.NotFound("config")
	}
	configURL := globalConfig.ConfigURL
	configMutex.Unlock()

	if configURL == "" {
		return nil, errors.NotFound("config")
	}

	log.Printf("[Worker] Executing task: GET %s", configURL)

	req, err := http.NewRequestWithContext(ctx, "GET", configURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to create request")
	}

	req.Header.Set("User-Agent", "curl/7.81.0")
	req.Header.Set("Accept", "text/plain")

	resp, err := u.httpClient.Do(req)
    if err != nil {
        return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to send request")
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to read response body")
    }

    var result any
    if err := json.Unmarshal(bodyBytes, &result); err != nil {
        log.Printf("[Worker] Task completed: Status %d, Non-JSON response", resp.StatusCode)
        return string(bodyBytes), nil
    }

    log.Printf("[Worker] Task completed: Status %d, JSON response", resp.StatusCode)
    return string(bodyBytes), nil
}

func (u *Worker) UpdateConfig(ctx context.Context, req worker.UpdateConfigRequest) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	globalConfig = &worker.WorkerConfig{
		ConfigURL:       req.ConfigURL,
		PoolingInterval: req.PoolingInterval,
		Version:         req.Version,
		UUID:            req.UUID,
	}

	log.Printf("============================================================")
	log.Println("[Worker] CONFIG UPDATED FROM AGENT!")
	log.Printf("  UUID: %s", globalConfig.UUID)
	log.Printf("  Version: %d", globalConfig.Version)
	log.Printf("  Config URL: %s", globalConfig.ConfigURL)
	log.Printf("  Pooling Interval: %d seconds", globalConfig.PoolingInterval)
	log.Printf("============================================================")

	return nil
}