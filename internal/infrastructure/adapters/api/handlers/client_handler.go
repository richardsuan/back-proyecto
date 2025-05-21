package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/richardsuan/back-proyecto/internal/core/ports/driver"
)

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
	var request struct {
		ClientName string `json:"client_name"`
	}

	// Lee el cuerpo JSON de la solicitud
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	clientName := request.ClientName
	log.Printf("Client name: %s", clientName)

	// Obtén los datos del cliente
	clientData, err := h.apiService.GetClientData(clientName)
	if err != nil {
		log.Printf("Error fetching client data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(clientData) == 0 {
		log.Printf("No data found for client: %s", clientName)
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}

	// Procesa los datos para asegurarte de que están en el formato correcto
	responseData := []map[string]interface{}{}
	for _, entry := range clientData {
		processedEntry := map[string]interface{}{
			"Fecha":       entry["Fecha"],
			"Presion":     entry["Presion"],
			"Temperatura": entry["Temperatura"],
			"Volumen":     entry["Volumen"],
			"Anomalia":    entry["Anomalia"], // Agregar el valor de anomalia
		}
		responseData = append(responseData, processedEntry)
	}

	log.Printf("Response data: %v", responseData)
	c.JSON(http.StatusOK, responseData)
}
