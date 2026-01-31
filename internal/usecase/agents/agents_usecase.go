package agents

import (
	"context"
	"distributed_system/internal/config"
	"distributed_system/internal/domain/agents"
	"distributed_system/pkg/crypto"
	"distributed_system/pkg/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AgentUsecase struct {
	repository agents.Repostiory
	cfg        *config.Config
}

func NewAgentUsecase(repository agents.Repostiory, cfg *config.Config) agents.Usecase {
	return &AgentUsecase{repository: repository, cfg: cfg}
}

func (u *AgentUsecase) Create(ctx context.Context) (string, error) {
	now := time.Now().Format(time.RFC3339)

	agent := &agents.Agent{
		UUID:        uuid.New().String(),
		CreatedAt: now,
	}

	if err := u.repository.Create(ctx, agent); err != nil {
		return "", errors.Wrap(err, "agent", "failed to create agent")
	}

	tokenAccessConfig, err := crypto.Generate(agent.UUID, u.cfg.Security.AgentSig)
	if err != nil {
		return "", errors.Wrap(err, "agent", "failed to create access token")
	}

	return tokenAccessConfig, nil
}

func (u *AgentUsecase) CreateRegistrationToken(ctx context.Context) (string, error) {
	token, err := bcrypt.GenerateFromPassword([]byte(u.cfg.Security.AgentSecret), 10)
	if err != nil {
		return "", errors.Wrap(err, "agent", "failed to create registration token")
	}

	return string(token), nil
}