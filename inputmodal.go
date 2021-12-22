package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalInput is based on Modal from tview, but has an input field instead
type ModalInput struct {
	*tview.Form
	frame *tview.Frame
	text  string
	done  func(string, bool)
}

func NewModalInput() *ModalInput {
	form := tview.NewForm()
	m := &ModalInput{form, tview.NewFrame(form), "", nil}

	m.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(0, 0, 0, 0)
	m.AddInputField("", "", 50, nil, func(text string) {
		m.text = text
	})
	m.AddButton("OK", func() {
		if m.done != nil {
			m.done(m.text, true) // Passed
		}
	})
	m.AddButton("Cancel", func() {
		if m.done != nil {
			m.done(m.text, false)
		}
	})
	m.frame.SetBorders(0, 0, 1, 0, 0, 0).
		SetBorder(true).
		SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)

	return m
}

// SetValue sets the current value in the item
func (m *ModalInput) SetValue(text string) {
	m.Clear(false)
	m.AddInputField("", text, 50, nil, func(text string) {
		m.text = text
	})
}

// SetDoneFunc sets the done func for this input.
// Will be called with the text of the input and a boolean for OK or cancel button.
func (m *ModalInput) SetDoneFunc(handler func(string, bool)) *ModalInput {
	m.done = handler
	return m
}

// Draw draws this primitive onto the screen.
func (m *ModalInput) Draw(screen tcell.Screen) {
	// Calculate the width of this modal.
	buttonsWidth := 50
	screenWidth, screenHeight := screen.Size()
	width := screenWidth / 3
	if width < buttonsWidth {
		width = buttonsWidth
	}
	// width is now without the box border.

	// Set the modal's position and size.
	height := 7
	width += 4
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}
