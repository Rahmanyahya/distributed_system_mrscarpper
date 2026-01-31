package handler

import (
	"distributed_system/internal/domain/admin"
	"distributed_system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	usecase admin.Usecase
}

func NewAdminHandler(usecase admin.Usecase) *AdminHandler {
	return &AdminHandler{usecase: usecase}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var input admin.InputLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BindingError(c, err)
		return
	}

	token, err := h.usecase.Login(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, token)
}