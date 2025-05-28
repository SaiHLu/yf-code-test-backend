package handler

import (
	docs "codetest/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerHandler struct {
	router *gin.RouterGroup
}

func NewSwaggerHandler(router *gin.RouterGroup) *SwaggerHandler {
	handler := &SwaggerHandler{
		router: router,
	}

	handler.registerRoutes()

	return handler
}

func (h *SwaggerHandler) registerRoutes() {
	docs.SwaggerInfo.BasePath = "/api"
	route := h.router.Group("/swagger")
	{
		route.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
