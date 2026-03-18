package handler

import (
	"github.com/VAGRAMCHIC/vless_reality_agent/internal/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userH *UserHandler
	log   *utils.Logger
}

func NewHandler(u *UserHandler, log *utils.Logger) *Handler {
	return &Handler{userH: u, log: log}
}

func (h *Handler) Register(r *gin.Engine, apiKey string) {
	r.Use(APIKeyMiddleware(apiKey))
	r.Use(RequestIDMiddleware(h.log))

	r.POST("/users", h.userH.AddUser)
	r.GET("/users", h.userH.ShowAllUsers)
	r.DELETE("/users/:id", h.userH.DelUser)
}
