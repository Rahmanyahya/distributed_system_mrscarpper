package handler

import (
	"context"
	"distributed_system/internal/domain/config"
	"distributed_system/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	config config.Usecase
}

func NewConfigHandler(config config.Usecase) *ConfigHandler {
	return &ConfigHandler{config: config}
}

func (h *ConfigHandler) GetLatestConfigAdmin(c *gin.Context) {
	config, err := h.config.GetLatestConfig(context.Background(), nil)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, config)
}

func (h *ConfigHandler) GetLatestConfigModel(c *gin.Context) {
	uuid, exist := c.Get("uuid")

	if !exist {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	uuidStr, ok := uuid.(string)
	if !ok {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	
	config, err := h.config.GetLatestConfig(context.Background(), &uuidStr)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, config)
}

func (h *ConfigHandler) Create(gin *gin.Context) {
	var input config.SaveCreate

	if err := gin.ShouldBindJSON(&input); err != nil {
		response.BindingError(gin, err)
		return
	}

	config, err := h.config.Create(context.Background(), &input)
	if err != nil {
		response.Error(gin, err)
		return
	}

	response.Success(gin, config)
}

func (h *ConfigHandler) Update(gin *gin.Context) {
	var input config.SaveUpdate

	if err := gin.ShouldBindJSON(&input); err != nil {
		response.BindingError(gin, err)
		return
	}

	if err := h.config.Update(context.Background(), &input); err != nil {
		response.Error(gin, err)
		return
	}

	response.Success(gin, nil)
}