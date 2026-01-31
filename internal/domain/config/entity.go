package config

import "context"

type Config struct {
	UUID      string `json:"uuid" gorm:"column:uuid;type:text;primaryKey"`
	Version   int  `json:"version" gorm:"column:version;type:int"`
	ConfigURL string `json:"config_url" gorm:"column:config_url;type:text"`
	PoolingInterval int `json:"pooling_interval" gorm:"column:pooling_interval;type:int"`
	CreatedAt string `json:"created_at" gorm:"column:created_at;type:text"`
}

func (Config) TableName() string {
	return "config"
}

type Repository interface {
	GetLatestConfig(ctx context.Context) (*Config, error)
	Create(ctx context.Context, config *Config) error
	Update(ctx context.Context, config *Config) error
}

type Usecase interface {
	GetLatestConfig(ctx context.Context, agentID *string) (*Config, error)
	Create(ctx context.Context, save *SaveCreate) (*Config, error)
	Update(ctx context.Context, save *SaveUpdate) error
}

type SaveCreate struct {
	ConfigUrl string `json:"config_url" binding:"required"`
	PoolingInterval int `json:"pooling_interval" binding:"min=30"`
}

type SaveUpdate struct {
	ConfigUrl string `json:"config_url" binding:"omitempty"`
	PoolingInterval *int `json:"pooling_interval" binding:"omitempty,min=30"`
}