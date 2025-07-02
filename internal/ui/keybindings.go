package ui

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
)

func Keybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.OnCursor < len(m.Connections)-1 {
			m.OnCursor++
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.OnCursor > 0 {
			m.OnCursor--
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Connections", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.Connections) > 0 {
			m.Selected = m.OnCursor
			tables, err := connection.GetAllTables()
			if err == nil {
				m.Tables = tables
				m.TablesCursor = 0
			}
			updateViews(g, m, connection)
			if _, err := g.SetCurrentView("Tables"); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Connections", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Tables"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Connections"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.TablesCursor < len(m.Tables)-1 {
			m.TablesCursor++
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.TablesCursor > 0 {
			m.TablesCursor--
			updateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.Tables) > 0 {
			tableName := m.Tables[m.TablesCursor]
			tableData, err := connection.GetTableData(tableName)
			if err == nil {
				m.TableData = populateTableData(tableData)
				m.DataCursorRow = 1
				m.DataCursorCol = 0
				updateViews(g, m, connection)
				if v, err := g.View("Data"); err == nil {
					v.SetOrigin(0, 0)
				}
				if _, err := g.SetCurrentView("Data"); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Data", 'q', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Tables"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Data", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorCol < len(m.TableData[0])-1 {
			colWidth := getColumnWidth(m.TableData, m.DataCursorCol)
			m.DataCursorCol++
			updateViews(g, m, connection)
			return scrollView(v, colWidth, 0)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorCol > 0 {
			m.DataCursorCol--
			colWidth := getColumnWidth(m.TableData, m.DataCursorCol)
			updateViews(g, m, connection)
			return scrollView(v, -colWidth, 0)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.DataCursorRow > 1 {
			m.DataCursorRow--
			updateViews(g, m, connection)
		}
		return scrollView(v, 0, -1)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorRow < len(m.TableData)-1 {
			m.DataCursorRow++
			updateViews(g, m, connection)

			return scrollView(v, 0, 1)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return ShowInputPopup(g, "Query", "Enter SQL Query", "", func(input string) error {
			if input != "" {
				tableName := m.Tables[m.TablesCursor]
				tableData, err := connection.GetTableData(tableName, input)
				if err == nil {
					m.TableData = populateTableData(tableData)
					m.DataCursorRow = 1
					m.DataCursorCol = 0

					updateViews(g, m, connection)
					if _, err := g.SetCurrentView("Data"); err != nil {
						return err
					}
				}
			}
			return nil
		}, func() error {
			g.SetCurrentView("Data")
			return nil
		})
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", 'y', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorRow < len(m.TableData) && m.DataCursorCol < len(m.TableData[m.DataCursorRow]) {
			cellValue := m.TableData[m.DataCursorRow][m.DataCursorCol]

			cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | pbcopy", cellValue))
			err := cmd.Run()
			if err != nil {
				return err
			}

			m.CellBlinking = true
			updateViews(g, m, connection)

			go func() {
				time.Sleep(200 * time.Millisecond)
				g.Update(func(g *gocui.Gui) error {
					m.CellBlinking = false
					updateViews(g, m, connection)
					return nil
				})
			}()
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Data", 'e', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorRow < len(m.TableData) && m.DataCursorCol < len(m.TableData[m.DataCursorRow]) {
			cellValue := m.TableData[m.DataCursorRow][m.DataCursorCol]
			columnName := m.TableData[0][m.DataCursorCol]
			tableName := m.Tables[m.TablesCursor]

			return ShowInputPopup(g, "Edit Value", fmt.Sprintf("Edit value for %s", columnName), cellValue, func(input string) error {
				if input != cellValue {
					primaryKey, err := connection.GetPrimaryKeyColumn(tableName)
					if err != nil {
						return err
					}

					pkIndex := -1
					for i, header := range m.TableData[0] {
						if header == primaryKey {
							pkIndex = i
							break
						}
					}

					if pkIndex == -1 {
						return fmt.Errorf("Primary key column not found")
					}

					pkValue := m.TableData[m.DataCursorRow][pkIndex]

					err = connection.UpdateTableValue(tableName, columnName, primaryKey, pkValue, input)
					if err != nil {
						return err
					}

					m.TableData[m.DataCursorRow][m.DataCursorCol] = input
					updateViews(g, m, connection)
				}
				return nil
			}, func() error {
				g.SetCurrentView("Data")
				return nil
			})
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Tables", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Data"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func scrollView(v *gocui.View, dx, dy int) error {
	if v != nil {
		ox, oy := v.Origin()
		if ox+dx >= 0 && oy+dy >= 0 {
			v.SetOrigin(ox+dx, oy+dy)
		}
	}
	return nil
}

func populateTableData(tableData *postgres.TableData) [][]string {
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

func getColumnWidth(tableData [][]string, colIndex int) int {
	if len(tableData) == 0 || colIndex < 0 || (len(tableData) > 0 && colIndex >= len(tableData[0])) {
		return 0
	}

	maxWidth := 0
	for _, row := range tableData {
		if colIndex < len(row) && len(row[colIndex]) > maxWidth {
			maxWidth = len(row[colIndex])
		}
	}

	// Add some padding (e.g., 2 spaces)
	return maxWidth + 2
}
