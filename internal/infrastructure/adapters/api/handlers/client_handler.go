package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/richardsuan/back-proyecto/internal/core/ports/driver"
)

type ClientDataRequest struct {
	ClientName string `json:"client_name"`
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

	// Obtén los datos de la hoja correspondiente al cliente
	rawData, err := h.apiService.GetClientRawData(clientName)
	if err != nil {
		log.Printf("Error fetching raw data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(rawData) < 2 {
		log.Printf("Not enough data: %v", rawData)
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}

	headers := rawData[0]
	log.Printf("Headers: %v", headers)

	responseData := []map[string]interface{}{}
	for _, row := range rawData[1:] {
		if len(row) != len(headers) {
			log.Printf("Skipping inconsistent row: %v", row)
			continue
		}

		entry := map[string]interface{}{}
		for i, header := range headers {
			value := row[i]
			if header == "Fecha" {
				fechaLimpia := strings.ReplaceAll(value, "\\", "")
				if t, err := time.Parse("2006-01-02 15:04:05", fechaLimpia); err == nil {
					// Convierte la fecha al formato ISO 8601
					entry[header] = t.Format("2006-01-02T15:04:05")
				} else {
					log.Printf("Error parsing fecha: '%s', error: %v", value, err)
					entry[header] = nil
				}
			} else {
				valorLimpio := strings.ReplaceAll(value, ",", ".")
				if num, err := strconv.ParseFloat(valorLimpio, 64); err == nil {
					entry[header] = num
				} else {
					log.Printf("Error parsing numeric value: '%s', error: %v", value, err)
					entry[header] = nil
				}
			}
		}
		responseData = append(responseData, entry)
	}

	log.Printf("Response data: %v", responseData)
	c.JSON(http.StatusOK, responseData)
}
