package application

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type PredictionService struct {
	modelPaths map[int]string // Map de cluster a modelo
	filePath   string         // Ruta del archivo CSV
	clusterMap map[string]int // Map de cliente a cluster
}

func NewPredictionService(modelPaths map[int]string, filePath string, clusterFilePath string) (*PredictionService, error) {
	// Cargar el mapa de cliente a cluster
	clusterMap, err := loadClusterMap(clusterFilePath)
	if err != nil {
		return nil, fmt.Errorf("error loading cluster map: %w", err)
	}

	return &PredictionService{
		modelPaths: modelPaths,
		filePath:   filePath,
		clusterMap: clusterMap,
	}, nil
}

func loadClusterMap(clusterFilePath string) (map[string]int, error) {
	file, err := os.Open(clusterFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	clusterMap := make(map[string]int)
	for _, row := range records[1:] { // Saltar la cabecera
		if len(row) < 2 {
			continue
		}
		cluster, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, fmt.Errorf("invalid cluster value: %v", err)
		}
		clusterMap[row[0]] = cluster
	}

	return clusterMap, nil
}

func (s *PredictionService) PredictAndSave(clientName string, data map[string]float64) (bool, error) {
	// Obtener el cluster del cliente
	cluster, ok := s.clusterMap[clientName]
	if !ok {
		return false, fmt.Errorf("client not found in cluster map: %s", clientName)
	}

	// Cargar el modelo correspondiente al cluster
	modelPath, ok := s.modelPaths[cluster]
	if !ok {
		return false, fmt.Errorf("model not found for cluster: %d", cluster)
	}
	fmt.Printf("Using model for cluster %d: %s\n", cluster, modelPath)

	// Preparar los datos para enviar al servicio REST
	requestBody := map[string]interface{}{
		"model_path":  modelPath,
		"presion":     data["Presion"],
		"temperatura": data["Temperatura"],
		"volumen":     data["Volumen"],
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return false, fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Enviar la solicitud al servicio REST
	resp, err := http.Post("http://localhost:5000/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta del servicio REST
	var response struct {
		IsAnomaly int `json:"is_anomaly"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, fmt.Errorf("error decoding response: %w", err)
	}

	// Determinar si es una anomalía
	isAnomaly := response.IsAnomaly == 1

	// Actualizar el archivo CSV
	err = s.updateCSV(data["Fecha"], isAnomaly)
	if err != nil {
		return false, fmt.Errorf("error updating CSV: %w", err)
	}

	return isAnomaly, nil
}

func (s *PredictionService) updateCSV(fecha float64, isAnomaly bool) error {
	file, err := os.OpenFile(s.filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Buscar la fila correspondiente y actualizar el campo "anomalia"
	for i, row := range records {
		if len(row) > 0 && row[0] == strconv.FormatFloat(fecha, 'f', -1, 64) {
			if isAnomaly {
				row[5] = "1" // Suponiendo que la columna "anomalia" está en el índice 5
			} else {
				row[5] = "0"
			}
			records[i] = row
			break
		}
	}

	// Reescribir el archivo CSV
	file.Truncate(0)
	file.Seek(0, 0)
	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}
