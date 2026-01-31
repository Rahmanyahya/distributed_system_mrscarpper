package agents

import "context"

type Agent struct {
	UUID        string `json:"uuid" gorm:"column:uuid;type:text;primaryKey"`
	CreatedAt string `json:"created_at" gorm:"column:created_at;type:text"`
}

func (Agent) TableName() string {
	return "agents"
}

type Repostiory interface {
	Create(ctx context.Context, agent *Agent) error
	GetById(ctx context.Context, ID string) (*Agent, error)
	GetAll(ctx context.Context) ([]Agent, error)
}

type Usecase interface {
	Create(ctx context.Context) (string, error)
	CreateRegistrationToken(ctx context.Context) (string, error)
}

