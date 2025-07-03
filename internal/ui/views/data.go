package views

import (
	"fmt"

	"github.com/rubengardner/lazy-database/backend/databases/postgres"
)

func PopulateTableData(tableData *postgres.TableData) [][]string {
	result := [][]string{}

	if tableData == nil {
		return result
	}

	if len(tableData.Headers) > 0 {
		result = append(result, tableData.Headers)
	}
	for _, row := range tableData.Rows {
		rowData := []string{}
		for _, val := range row {
			rowData = append(rowData, fmt.Sprintf("%v", val))
		}
		result = append(result, rowData)
	}

	return result
}
