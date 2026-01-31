package config

import (
	"context"
	"distributed_system/internal/domain/agents"
	"distributed_system/internal/domain/config"
	"distributed_system/internal/infrastructure/cache"
	"distributed_system/pkg/errors"
	"fmt"
	"time"

	configEnv "distributed_system/internal/config"

	"github.com/google/uuid"
)

type ConfigUsecase struct {
	repository config.Repository
	agentsRepository agents.Repostiory
	cfg        *configEnv.Config
	cache      *cache.ConfigCache
}

func NewConfigUsecase(repository config.Repository, agentRespository agents.Repostiory, cfg *configEnv.Config, cache *cache.ConfigCache) config.Usecase {
	return &ConfigUsecase{
		repository: repository,
		agentsRepository: agentRespository,
		cfg: cfg,
		cache: cache,
	}
}

func (u *ConfigUsecase) GetLatestConfig(ctx context.Context, agentID *string) (*config.Config, error) {
	if agentID != nil {
		_, err := u.agentsRepository.GetById(ctx, *agentID)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil, errors.NotFound("agent")
			}

			return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to get agent")
		}
	}

	chaced, err := u.cache.GetConfig(ctx)
	if err == nil && chaced != nil {
		return chaced, nil
	}

	config, err := u.repository.GetLatestConfig(ctx)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFound("config")
		}
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to get latest config")
	}

	u.cache.SetConfig(ctx, config)

	return config, nil
}

func (u *ConfigUsecase) Create(ctx context.Context, save *config.SaveCreate) (*config.Config, error) {
	now := time.Now().Format(time.RFC3339)

	var version int

	latestConfig, err := u.repository.GetLatestConfig(ctx)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	if errors.IsNotFound(err) {
		version = 1
	} else {
		version = latestConfig.Version + 1
	}

	newConfig := &config.Config{
		UUID:      uuid.New().String(),
		Version:   version,
		ConfigURL: save.ConfigUrl,
		PoolingInterval: save.PoolingInterval,
		CreatedAt: now,
	}

	if err := u.repository.Create(ctx, newConfig); err != nil {
		return nil, errors.Wrap(err, "config", "failed to create config")
	}

	if err := u.cache.SetConfig(ctx, newConfig); err != nil {
		return nil, errors.Wrap(err, "config", "failed to cache config")
	}

	return newConfig, nil
}

func (u *ConfigUsecase) Update(ctx context.Context, save *config.SaveUpdate) error {
	config, err := u.repository.GetLatestConfig(ctx)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.NotFound("config")
		}
		return errors.Wrap(err, "config", "failed to get config")
	}

	if save.ConfigUrl != "" {
		config.ConfigURL = save.ConfigUrl
	}

	if save.PoolingInterval != nil && *save.PoolingInterval >= 0 {
		config.PoolingInterval = *save.PoolingInterval
	}

	if err = u.repository.Update(ctx, config); err != nil {
		return errors.Wrap(err, "config", "failed to update config")
	}

	if err = u.cache.SetConfig(ctx, config); err != nil {
		fmt.Print(err)
		return errors.Wrap(err, "config", "failed to cache config")
	}

	return nil
}
