package application

import (
	"github.com/richardsuan/back-proyecto/internal/core/ports/driven"
	"github.com/richardsuan/back-proyecto/internal/core/ports/driver"
)

type ClientService struct {
	excelRepo driven.ExcelRepository
	filePath  string
}

func NewClientService(excelRepo driven.ExcelRepository, filePath string) driver.APIService {
	return &ClientService{
		excelRepo: excelRepo,
		filePath:  filePath,
	}
}

func (s *ClientService) GetClientNames() ([]string, error) {
	return s.excelRepo.GetSheetNames(s.filePath)
}

func (s *ClientService) GetClientRawData(clientName string) ([][]string, error) {
	return s.excelRepo.GetClientData(s.filePath, clientName)
}
