package excel

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/richardsuan/back-proyecto/internal/core/ports/driven"
)

type ExcelAdapter struct{}

func NewExcelAdapter() driven.ExcelRepository {
	return &ExcelAdapter{}
}

// GetSheetNames now retrieves distinct client names from the CSV file.
func (a *ExcelAdapter) GetSheetNames(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %w", err)
	}

	clientSet := make(map[string]struct{})
	for _, row := range records[1:] { // Skip the header row
		if len(row) < 5 {
			continue
		}
		clientName := strings.TrimSpace(row[4]) // Column "Cliente" is at index 4
		if clientName != "" {
			clientSet[clientName] = struct{}{}
		}
	}

	var clientNames []string
	for client := range clientSet {
		clientNames = append(clientNames, client)
	}

	return clientNames, nil
}

// GetClientData retrieves all rows for a specific client from the CSV file.
func (a *ExcelAdapter) GetClientData(filePath string, clientName string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %w", err)
	}

	var data []map[string]interface{}
	for _, row := range records[1:] { // Skip the header row
		if len(row) < 7 { // Ensure the row has all required columns
			continue
		}
		if strings.TrimSpace(row[4]) == clientName { // Column "Cliente" is at index 4
			entry := map[string]interface{}{
				"Fecha":       strings.TrimSpace(row[0]), // Fecha
				"Presion":     parseFloat(row[1]),        // Presion
				"Temperatura": parseFloat(row[2]),        // Temperatura
				"Volumen":     parseFloat(row[3]),        // Volumen
				"Anomalia":    parseInt(row[5]),          // Anomalia
			}
			data = append(data, entry)
		}
	}

	return data, nil
}

// Helper function to parse a string to a float64
func parseFloat(value string) float64 {
	parsedValue, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0 // Default to 0 if parsing fails
	}
	return parsedValue
}

// Helper function to parse a string to an int
func parseInt(value string) int {
	parsedValue, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0 // Default to 0 if parsing fails
	}
	return parsedValue
}
