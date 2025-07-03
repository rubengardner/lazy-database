package data

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui"
	"github.com/rubengardner/lazy-database/internal/ui/views/views"
)

func QueryKeybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("Data", 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return ui.ShowInputPopup(g, "Query", "Enter SQL Query", "", func(input string) error {
			if input != "" {
				tableName := m.Tables[m.TablesCursor]
				tableData, err := connection.GetTableData(tableName, input)
				if err == nil {
					m.TableData = views.PopulateTableData(tableData)
					m.DataCursorRow = 1
					m.DataCursorCol = 0

					ui.UpdateViews(g, m, connection)
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

	return nil
}
