package ui

import "strings"

func calculateColumnWidths(data [][]string) []int {
	if len(data) == 0 {
		return []int{}
	}

	widths := make([]int, len(data[0]))
	for i, col := range data[0] {
		widths[i] = len(col)
	}

	for _, row := range data {
		for i, col := range row {
			if i < len(widths) && len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	return widths
}

func formatRowWithWidth(row []string, colWidths []int) string {
	var formattedRow string
	for i, cell := range row {
		if i < len(colWidths) {
			formattedCell := cell
			if len(cell) < colWidths[i] {
				formattedCell = cell + strings.Repeat(" ", colWidths[i]-len(cell))
			}
			formattedRow += formattedCell + "  "
		}
	}
	return formattedRow
}
