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

func NewGinAdapter(apiService driver.APIService, predictionService driver.PredictionService) *GinAdapter {
	engine := gin.Default()
	clientHandler := handlers.NewClientHandler(apiService)
	predictionHandler := handlers.NewPredictionHandler(predictionService)

	engine.GET("/clients", clientHandler.GetClientNames)
	engine.POST("/clients/data", clientHandler.GetClientData)
	engine.POST("/predict", predictionHandler.Predict) // Nuevo endpoint

	return &GinAdapter{
		engine:     engine,
		apiService: apiService,
	}
}

func (a *GinAdapter) Run(addr string) error {
	return a.engine.Run(addr)
}
