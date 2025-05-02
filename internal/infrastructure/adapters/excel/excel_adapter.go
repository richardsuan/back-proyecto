package excel

import (
	"fmt"
	"strings"

	"github.com/richardsuan/back-proyecto/internal/core/ports/driven"

	"github.com/tealeg/xlsx"
)

type ExcelAdapter struct{}

func NewExcelAdapter() driven.ExcelRepository {
	return &ExcelAdapter{}
}

func (a *ExcelAdapter) GetSheetNames(filePath string) ([]string, error) {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening excel file: %w", err)
	}

	var sheetNames []string
	for _, sheet := range xlFile.Sheets {
		sheetNames = append(sheetNames, sheet.Name)
	}
	return sheetNames, nil
}

func (a *ExcelAdapter) GetClientData(filePath string, sheetName string) ([][]string, error) {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening excel file: %w", err)
	}

	sheet, ok := xlFile.Sheet[sheetName]
	if !ok {
		return nil, fmt.Errorf("sheet '%s' not found", sheetName)
	}

	var data [][]string

	for _, row := range sheet.Rows {
		if row == nil {
			continue
		}
		var rowData []string
		for _, cell := range row.Cells {
			value := strings.TrimSpace(cell.String())
			rowData = append(rowData, value)
		}
		data = append(data, rowData)
	}

	return data, nil
}
