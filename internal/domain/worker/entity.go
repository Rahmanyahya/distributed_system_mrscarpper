package worker

import "context"

type WorkerConfig struct {
	ConfigURL       string `json:"config_url"`
	PoolingInterval int    `json:"pooling_interval"`
	Version         int    `json:"version"`
	UUID            string `json:"uuid"`
}

type UpdateConfigRequest struct {
	ConfigURL       string `json:"config_url" binding:"required"`
	PoolingInterval int    `json:"pooling_interval" binding:"required,min=30"`
	Version         int    `json:"version" binding:"required"`
	UUID            string `json:"uuid" binding:"required"`
}

type Usecase interface {
	Hit(ctx context.Context) (any, error)
	UpdateConfig(ctx context.Context, req UpdateConfigRequest) error
}