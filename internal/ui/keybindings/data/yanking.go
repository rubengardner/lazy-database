package data

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/rubengardner/lazy-database/backend/databases/postgres"
	"github.com/rubengardner/lazy-database/internal/model"
	"github.com/rubengardner/lazy-database/internal/ui"
)

func YankKeybinding(g *gocui.Gui, m *model.LazyDBState, connection *postgres.DatabaseConnection) error {
	if err := g.SetKeybinding("Data", 'y', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if len(m.TableData) > 0 && m.DataCursorRow < len(m.TableData) && m.DataCursorCol < len(m.TableData[m.DataCursorRow]) {
			cellValue := m.TableData[m.DataCursorRow][m.DataCursorCol]

			cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | pbcopy", cellValue))
			err := cmd.Run()
			if err != nil {
				return err
			}
			m.CellBlinking = true
			ui.UpdateViews(g, m, connection)

			go func() {
				time.Sleep(200 * time.Millisecond)
				g.Update(func(g *gocui.Gui) error {
					m.CellBlinking = false
					ui.UpdateViews(g, m, connection)
					return nil
				})
			}()
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
