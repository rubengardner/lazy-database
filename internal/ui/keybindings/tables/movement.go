package tables

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui"
	"github.com/rubengardner/lazy-database/internal/ui/views/views"
)

func MovementTablesKeybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("Tables", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.TablesCursor < len(m.Tables)-1 {
			m.TablesCursor++
			ui.UpdateViews(g, m, connection)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Tables", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.TablesCursor > 0 {
			m.TablesCursor--
			ui.UpdateViews(g, m, connection)
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
				m.TableData = views.PopulateTableData(tableData)
				m.DataCursorRow = 1
				m.DataCursorCol = 0
				ui.UpdateViews(g, m, connection)
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

	return nil
}
