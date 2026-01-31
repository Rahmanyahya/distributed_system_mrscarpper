package config

import (
	"context"
	"distributed_system/internal/domain/config"
	"distributed_system/internal/infrastructure/cache"
	"distributed_system/pkg/errors"

	"gorm.io/gorm"
)

type repository struct {
	db    *gorm.DB
	cache *cache.ConfigCache
}

func NewCOnfigRepository(db *gorm.DB, cache *cache.ConfigCache) config.Repository {
	return &repository{
		db:    db,
		cache: cache,
	}
}

func (r *repository) GetLatestConfig(ctx context.Context) (*config.Config, error) {
    var cfg config.Config
    res := r.db.WithContext(ctx).
        Order("version DESC").
        First(&cfg) // First otomatis menambahkan LIMIT 1
    
    if res.Error != nil {
        if errors.Is(res.Error, gorm.ErrRecordNotFound) {
            return nil, errors.NotFound("config")
        }
        return nil, errors.Database(res.Error)
    }

    return &cfg, nil
}


func (r *repository) Create(ctx context.Context, config *config.Config) error {
	err := r.db.WithContext(ctx).Create(config).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, config *config.Config) error {
	err := r.db.WithContext(ctx).Save(config).Error
	if err != nil {
		return err
	}

	return nil
}