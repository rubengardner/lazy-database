package ui

import (
	"fmt"

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
				m.TableData = [][]string{}
				m.TableData = [][]string{}

				if tableData != nil {
					headers := tableData.Headers
					if len(headers) > 0 {
						m.TableData = append(m.TableData, headers)
					}

					for _, row := range tableData.Rows {
						rowData := []string{}
						for _, val := range row {
							rowData = append(rowData, fmt.Sprintf("%v", val))
						}
						m.TableData = append(m.TableData, rowData)
					}
				}
			}
			updateViews(g, m, connection)
			if _, err := g.SetCurrentView("Data"); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Data", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if _, err := g.SetCurrentView("Tables"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("Data", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return scrollView(v, 2, 0) // Scroll right
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return scrollView(v, 0, -1) // Scroll up
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("Data", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return scrollView(v, 0, 1) // Scroll down
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
