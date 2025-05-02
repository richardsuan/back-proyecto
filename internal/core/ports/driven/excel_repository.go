package driven

type ExcelRepository interface {
	GetSheetNames(filePath string) ([]string, error)
	GetClientData(filePath string, sheetName string) ([][]string, error)
}
