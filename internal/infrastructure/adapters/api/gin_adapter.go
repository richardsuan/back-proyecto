package api

import (
	"github.com/gin-gonic/gin"
	"github.com/richardsuan/back-proyecto/internal/core/ports/driver"
	"github.com/richardsuan/back-proyecto/internal/infrastructure/adapters/api/handlers"
)

type GinAdapter struct {
	engine     *gin.Engine
	apiService driver.APIService
}

func NewGinAdapter(apiService driver.APIService) *GinAdapter {
	engine := gin.Default()
	handler := handlers.NewClientHandler(apiService)

	engine.GET("/clients", handler.GetClientNames)
	engine.POST("/clients/data", handler.GetClientData) // Nuevo endpoint

	return &GinAdapter{
		engine:     engine,
		apiService: apiService,
	}
}

func (a *GinAdapter) Run(addr string) error {
	return a.engine.Run(addr)
}
