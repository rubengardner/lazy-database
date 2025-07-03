package data

import (
	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui"
)

func MovementDataKeybindings(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("Data", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorCol < len(m.TableData[0])-1 {
			colWidth := getColumnWidth(m.TableData, m.DataCursorCol)
			m.DataCursorCol++
			ui.UpdateViews(g, m, connection)
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
			ui.UpdateViews(g, m, connection)
			return scrollView(v, -colWidth, 0)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if m.DataCursorRow > 1 {
			m.DataCursorRow--
			ui.UpdateViews(g, m, connection)
		}
		return scrollView(v, 0, -1)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorRow < len(m.TableData)-1 {
			m.DataCursorRow++
			ui.UpdateViews(g, m, connection)

			return scrollView(v, 0, 1)
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

	return nil
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
	return maxWidth + 2
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
