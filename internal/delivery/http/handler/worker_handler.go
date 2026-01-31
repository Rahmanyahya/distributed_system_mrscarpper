package handler

import (
	"distributed_system/internal/domain/worker"
	"distributed_system/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	usecase     worker.Usecase
	internalKey string
}

func NewWorkerHandler(usecase worker.Usecase) *WorkerHandler {
	return &WorkerHandler{
		usecase:     usecase,
		internalKey: "worker_internal_key_2024",
	}
}

func (h *WorkerHandler) Hit(c *gin.Context) {
	resp, err := h.usecase.Hit(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}

func (h *WorkerHandler) UpdateConfig(c *gin.Context) {
	var req worker.UpdateConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BindingError(c, err)
		return
	}

	log.Printf("[Worker] Received config update from Agent: Version=%d, URL=%s",
		req.Version, req.ConfigURL)

	if err := h.usecase.UpdateConfig(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Configuration updated successfully",
		"version": req.Version,
	})
}
