package main

import (
	"log"

	"github.com/richardsuan/back-proyecto/internal/core/application"
	"github.com/richardsuan/back-proyecto/internal/infrastructure/adapters/api"
	"github.com/richardsuan/back-proyecto/internal/infrastructure/adapters/excel"
)

func main() {
	excelAdapter := excel.NewExcelAdapter()
	clientService := application.NewClientService(excelAdapter, "./Datos.xlsx") // Replace with your file path
	apiAdapter := api.NewGinAdapter(clientService)

	err := apiAdapter.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to run API server: %v", err)
	}
}
