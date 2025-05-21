package main

import (
	"log"

	"github.com/richardsuan/back-proyecto/internal/core/application"
	"github.com/richardsuan/back-proyecto/internal/infrastructure/adapters/api"
	"github.com/richardsuan/back-proyecto/internal/infrastructure/adapters/excel"
)

func main() {
	// Definir las rutas de los modelos
	modelPaths := map[int]string{
		0: "./models/modelo_isolation_cluster_0.pkl",
		1: "./models/modelo_isolation_cluster_1.pkl",
		2: "./models/modelo_isolation_cluster_2.pkl",
		3: "./models/modelo_isolation_cluster_3.pkl",
	}

	// Inicializar el servicio de predicción
	predictionService, err := application.NewPredictionService(modelPaths, "./data/Data_final.csv", "./data/Cluster.csv")
	if err != nil {
		log.Fatalf("Error initializing prediction service: %v", err)
	}

	// Inicializar el adaptador de Excel
	excelAdapter := excel.NewExcelAdapter()

	// Inicializar el servicio de cliente
	clientService := application.NewClientService(excelAdapter, "./data/Data_final.csv")

	// Inicializar el adaptador de API
	apiAdapter := api.NewGinAdapter(clientService, predictionService)

	// Ejecutar el servidor
	err = apiAdapter.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
