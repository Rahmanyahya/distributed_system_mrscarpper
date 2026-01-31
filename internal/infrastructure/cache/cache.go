package cache

import (
	"context"
	"distributed_system/internal/domain/config"
	"distributed_system/internal/infrastructure/redis"
	"encoding/json"
	"time"
)

const (
	LatestConfigKey    = "config:latest"  
	DefaultCacheTTL = 24 * time.Hour * 30
)


type ConfigCache struct {
	redis *redis.Client
}

func NewConfigCache(redisClient *redis.Client) *ConfigCache {
	return &ConfigCache{
		redis: redisClient,
	}
}

func (c *ConfigCache) SetConfig(ctx context.Context, config *config.Config) error {
	data, err :=  json.Marshal(config)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, LatestConfigKey, data, DefaultCacheTTL)
}

func (c *ConfigCache) GetConfig(ctx context.Context) (*config.Config, error) {
	var cfg config.Config

	value, err := c.redis.Get(ctx, LatestConfigKey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(value), &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

