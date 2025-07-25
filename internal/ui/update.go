package ui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
)

func updateConnectionsView(v *gocui.View, m *model.LazyDBState) {
	v.Clear()
	for i, db := range m.Connections {
		cursor := " "
		if i == m.OnCursor {
			cursor = "❯"
		}
		fmt.Fprintf(v, "%s %s\n", cursor, db)
	}
}

func updateTablesView(v *gocui.View, m *model.LazyDBState, connection *postgres.DatabaseConnection) {
	v.Clear()
	if len(m.Connections) == 0 {
		fmt.Fprintf(v, "No database connections available")
		return
	}

	if len(m.Tables) == 0 {
		tables, err := connection.GetAllTables()
		if err != nil {
			fmt.Fprintf(v, "Error fetching tables: %v\n", err)
			return
		}
		m.Tables = tables
	}

	for i, table := range m.Tables {
		cursor := " "
		if i == m.TablesCursor {
			cursor = "❯"
		}
		fmt.Fprintf(v, "%s %s\n", cursor, table)
	}
}

func updateDataView(v *gocui.View, m *model.LazyDBState) {
	v.Clear()

	if len(m.TableData) == 0 {
		fmt.Fprintf(v, "Select a table to view data")
		return
	}

	colWidths := calculateColumnWidths(m.TableData)

	if len(m.TableData) > 0 {
		header := formatRowWithWidth(m.TableData[0], colWidths)
		fmt.Fprintf(v, "\033[1;37m%s\033[0m\n", header)

		separator := ""
		for _, width := range colWidths {
			separator += strings.Repeat("-", width) + "  "
		}
		fmt.Fprintf(v, "%s\n", separator)

		for i := 1; i < len(m.TableData); i++ {
			row := m.TableData[i]
			rowStr := ""

			for j, cell := range row {
				if i == m.DataCursorRow && j == m.DataCursorCol {
					if m.CellBlinking {
						// Blinking effect (using a different color or formatting)
						rowStr += fmt.Sprintf("\033[7;32m%-*s\033[0m  ", colWidths[j], cell)
					} else {
						// Normal selected cell highlighting
						rowStr += fmt.Sprintf("\033[7m%-*s\033[0m  ", colWidths[j], cell)
					}
				} else {
					rowStr += fmt.Sprintf("%-*s  ", colWidths[j], cell)
				}
			}

			fmt.Fprintf(v, "%s\n", rowStr)
		}
	}
}

func UpdateViews(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) {
	if v, err := g.View("Connections"); err == nil {
		updateConnectionsView(v, m)
	}
	if v, err := g.View("Tables"); err == nil {
		updateTablesView(v, m, connection)
	}
	if v, err := g.View("Data"); err == nil {
		updateDataView(v, m)
	}
}
