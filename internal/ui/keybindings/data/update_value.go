package data

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui"
)

func UpdateValueDataKeybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("Data", 'e', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorRow < len(m.TableData) && m.DataCursorCol < len(m.TableData[m.DataCursorRow]) {
			cellValue := m.TableData[m.DataCursorRow][m.DataCursorCol]
			columnName := m.TableData[0][m.DataCursorCol]
			tableName := m.Tables[m.TablesCursor]

			return ui.ShowInputPopup(g, "Edit Value", fmt.Sprintf("Edit value for %s", columnName), cellValue, func(input string) error {
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
						return fmt.Errorf("primary key column not found")
					}

					pkValue := m.TableData[m.DataCursorRow][pkIndex]

					err = connection.UpdateTableValue(tableName, columnName, primaryKey, pkValue, input)
					if err != nil {
						return err
					}

					m.TableData[m.DataCursorRow][m.DataCursorCol] = input
					ui.UpdateViews(g, m, connection)
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

	return nil
}
