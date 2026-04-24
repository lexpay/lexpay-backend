package auth

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Svc
}

func NewHandler(service Svc) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleSignUp(c *gin.Context) {}