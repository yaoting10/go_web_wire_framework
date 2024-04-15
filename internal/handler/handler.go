package handler

import (
	"github.com/gin-gonic/gin"
	"goboot/pkg/log"
)

type Handler struct {
	L *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		L: logger,
	}
}

func (h *Handler) Group(path string, e *gin.Engine, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return e.Group(path, handlers...)
}
