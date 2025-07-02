package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type PopupOptions struct {
	Title      string
	Message    string
	Width      int
	Height     int
	Centered   bool
	OnClose    func() error
	InputField bool
	OnSubmit   func(string) error
}

// CreatePopup creates and displays a popup in the UI
func CreatePopup(g *gocui.Gui, opts PopupOptions) error {
	maxX, maxY := g.Size()

	width := opts.Width
	if width <= 0 {
		width = 40 // default width
	}

	height := opts.Height
	if height <= 0 {
		height = 6 // default height
	}

	// Calculate position for centering if requested
	x, y := 0, 0
	if opts.Centered {
		x = (maxX - width) / 2
		y = (maxY - height) / 2
	}

	return createPopupView(g, x, y, width, height, opts)
}

func createPopupView(g *gocui.Gui, x, y, width, height int, opts PopupOptions) error {
	// Create popup view
	popupName := "popup"
	if v, err := g.SetView(popupName, x, y, x+width, y+height); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = opts.Title
		v.Wrap = true
		v.Editable = false
		fmt.Fprintln(v, opts.Message)

		if opts.InputField {
			inputName := "popupInput"
			inputWidth := width - 4

			if inputView, err := g.SetView(inputName, x+2, y+height-3, x+2+inputWidth, y+height-1); err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}

				inputView.Editable = true
				inputView.Wrap = true
				inputView.Highlight = true
				inputView.SelBgColor = gocui.ColorWhite
				inputView.SelFgColor = gocui.ColorBlack
				g.SetCurrentView(inputName)
				if err := g.SetKeybinding(inputName, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
					input := v.Buffer()
					g.DeleteView(inputName)
					g.DeleteView(popupName)

					if opts.OnSubmit != nil {
						return opts.OnSubmit(input)
					}
					return nil
				}); err != nil {
					return err
				}

				if err := g.SetKeybinding(inputName, 'q', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
					g.DeleteView(inputName)
					g.DeleteView(popupName)

					if opts.OnClose != nil {
						return opts.OnClose()
					}
					return nil
				}); err != nil {
					return err
				}
			}
		}

		if err := g.SetKeybinding(popupName, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			if opts.InputField {
				g.DeleteView("popupInput")
			}
			g.DeleteView(popupName)

			if opts.OnClose != nil {
				return opts.OnClose()
			}
			return nil
		}); err != nil {
			return err
		}

		if !opts.InputField {
			g.SetCurrentView(popupName)
		}
	}

	return nil
}

func ShowInfoPopup(g *gocui.Gui, title, message string) error {
	return CreatePopup(g, PopupOptions{
		Title:   title,
		Message: message,
		Width:   len(message) + 10, Height: 6,
		Centered: true,
	})
}

func ShowErrorPopup(g *gocui.Gui, message string) error {
	return CreatePopup(g, PopupOptions{
		Title:    "Error",
		Message:  message,
		Width:    len(message) + 10,
		Height:   6,
		Centered: true,
	})
}

func ShowConfirmPopup(g *gocui.Gui, message string, onConfirm func() error) error {
	confirmMessage := message + "\n\nPress Enter to confirm or Esc to cancel."

	return CreatePopup(g, PopupOptions{
		Title:    "Confirm",
		Message:  confirmMessage,
		Width:    len(confirmMessage) + 6,
		Height:   8,
		Centered: true,
		OnSubmit: func(_ string) error {
			if onConfirm != nil {
				return onConfirm()
			}
			return nil
		},
	})
}

func ShowInputPopup(g *gocui.Gui, title, prompt string, initialValue string, onSubmit func(string) error, onClose func() error) error {
	opts := PopupOptions{
		Title:      title,
		Message:    prompt,
		Width:      50,
		Height:     8,
		Centered:   true,
		InputField: true,
		OnSubmit:   onSubmit,
		OnClose:    onClose,
	}

	err := CreatePopup(g, opts)
	if err != nil {
		return err
	}

	if initialValue != "" {
		if v, err := g.View("popupInput"); err == nil {
			v.Clear()
			fmt.Fprint(v, initialValue)
			v.SetCursor(len(initialValue), 0)
		}
	}

	return nil
}
