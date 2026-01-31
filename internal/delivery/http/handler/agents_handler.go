package handler

import (
	"distributed_system/internal/domain/agents"
	"distributed_system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AgentsHandler struct {
	agentUsecase agents.Usecase
}

func NewAgentsHandler(agentUsecase agents.Usecase) *AgentsHandler {
	return &AgentsHandler{agentUsecase: agentUsecase}
}

func (h *AgentsHandler) Register(c *gin.Context) {
	token, err := h.agentUsecase.Create(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, token)
}

func (h *AgentsHandler) GenerateRegistrationConfifg(c *gin.Context) {
	token, err := h.agentUsecase.CreateRegistrationToken(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, token)
}