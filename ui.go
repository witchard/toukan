package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	col1 := tview.NewList()
	col1.ShowSecondaryText(false).SetTitle("To Do").SetBorder(true)
	col2 := tview.NewList()
	col2.ShowSecondaryText(false).SetTitle("Doing").SetBorder(true)
	col3 := tview.NewList()
	col3.ShowSecondaryText(false).SetTitle("Done").SetBorder(true)

	col1.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			return nil
		case tcell.KeyRight:
			return nil
		}
		switch event.Rune() {
		case 'q':
			app.Stop()
		}
		return event
	})

	for i := 0; i < 10; i++ {
		col1.AddItem(fmt.Sprint("Item ", i), "", 0, nil)
		col2.AddItem(fmt.Sprint("Item ", i), "", 0, nil)
		col3.AddItem(fmt.Sprint("Item ", i), "", 0, nil)
	}

	flex := tview.NewFlex().
		AddItem(col1, 0, 1, true).
		AddItem(col2, 0, 1, false).
		AddItem(col3, 0, 1, false)

	pages := tview.NewPages().
		AddPage("TouKan", flex, true, true)
	app.SetRoot(pages, true)

	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}
