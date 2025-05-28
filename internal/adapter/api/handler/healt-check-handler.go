package handler

import (
	"codetest/internal/adapter/api/presenter"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	router *gin.RouterGroup
}

func NewHealthCheckHandler(router *gin.RouterGroup) *HealthCheckHandler {
	handler := &HealthCheckHandler{
		router: router,
	}

	handler.registerRoutes()

	return handler
}

func (h *HealthCheckHandler) registerRoutes() {
	route := h.router.Group("/health")
	{
		route.GET("", h.HealthCheck)
	}
}

// HealtCheck godoc
// @Summary Health Check
// @Description Check the health of the service
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} presenter.JsonResponse
// @Router /health [get]
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, presenter.JsonResponse{
		Success: true,
		Data: map[string]string{
			"status":  "ok",
			"message": "Service is running",
		},
	})
}
