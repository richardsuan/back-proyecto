package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/richardsuan/back-proyecto/internal/core/ports/driver"
)

type PredictionHandler struct {
	predictionService driver.PredictionService
}

func NewPredictionHandler(predictionService driver.PredictionService) *PredictionHandler {
	return &PredictionHandler{
		predictionService: predictionService,
	}
}

func (h *PredictionHandler) Predict(c *gin.Context) {
	var request struct {
		ClientName  string  `json:"client_name"`
		Fecha       float64 `json:"fecha"`
		Presion     float64 `json:"presion"`
		Temperatura float64 `json:"temperatura"`
		Volumen     float64 `json:"volumen"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	data := map[string]float64{
		"Fecha":       request.Fecha,
		"Presion":     request.Presion,
		"Temperatura": request.Temperatura,
		"Volumen":     request.Volumen,
	}

	isAnomaly, err := h.predictionService.PredictAndSave(request.ClientName, data)
	if err != nil {
		log.Printf("Error predicting anomaly: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"anomalia": isAnomaly})
}

type PredictionService interface {
	PredictAndSave(clientName string, data map[string]float64) (bool, error)
}
