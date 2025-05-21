package driven

type ExcelRepository interface {
	GetSheetNames(filePath string) ([]string, error)
	GetClientData(filePath string, clientName string) ([]map[string]interface{}, error) // Cambiado a []map[string]interface{}
}
