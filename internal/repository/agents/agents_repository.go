package agents

import (
	"context"
	"distributed_system/internal/domain/agents"
	"distributed_system/pkg/errors"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) agents.Repostiory {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, agent *agents.Agent) error {
	return r.db.WithContext(ctx).Create(agent).Error
}

func (r *repository) GetById(ctx context.Context, ID string) (*agents.Agent, error) {
	var agent agents.Agent
	if err := r.db.First(&agent, "uuid = ?", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("agent")
		}
		return nil, errors.Database(err)
	}

	return &agent, nil
}

func (r *repository) GetAll(ctx context.Context) ([]agents.Agent, error) {
	var agents []agents.Agent
	if err := r.db.Find(&agents).Error; err != nil {
		return nil, errors.Database(err)
	}

	return agents, nil
}