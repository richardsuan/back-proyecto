package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richardsuan/back-proyecto/internal/core/ports/driver"
)

type ClientDataRequest struct {
	ClientName string `json:"client_name"`
	Variable   string `json:"variable"` // "Volumen", "Presion", "Temperatura"
}

type ClientHandler struct {
	apiService driver.APIService
}

func NewClientHandler(apiService driver.APIService) *ClientHandler {
	return &ClientHandler{
		apiService: apiService,
	}
}

func (h *ClientHandler) GetClientNames(c *gin.Context) {
	names, err := h.apiService.GetClientNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, names)
}

func (h *ClientHandler) GetClientData(c *gin.Context) {
	var request ClientDataRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	clientName := request.ClientName
	variable := request.Variable

	rawData, err := h.apiService.GetClientRawData(clientName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rawData) < 2 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}

	headers := rawData[0]
	fechaIndex := -1
	variableIndex := -1

	for i, h := range headers {
		if h == "Fecha" {
			fechaIndex = i
		}
		if h == variable {
			variableIndex = i
		}
	}

	if fechaIndex == -1 || variableIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Required columns not found"})
		return
	}

	responseData := []map[string]interface{}{}
	for _, row := range rawData[1:] {
		if len(row) <= variableIndex || len(row) <= fechaIndex {
			continue
		}
		fechaOriginal := row[fechaIndex]
		var timestamp int64 = 0

		// Limpia los caracteres de escape de la fecha
		fechaLimpia := strings.ReplaceAll(fechaOriginal, "\\", "")

		// Intenta parsear la fecha con el formato esperado
		if t, err := time.Parse("2006-01-02 15:04:05", fechaLimpia); err == nil {
			timestamp = t.Unix()
			log.Printf("Fecha parseada correctamente: timestamp = %d", timestamp)
		} else {
			log.Printf("Error parsing fecha: '%s', error: %v", fechaOriginal, err)
			timestamp = 0 // O algún valor por defecto
		}

		entry := map[string]interface{}{
			"Fecha":  timestamp,
			variable: row[variableIndex],
		}
		responseData = append(responseData, entry)
	}

	c.JSON(http.StatusOK, responseData)
}
