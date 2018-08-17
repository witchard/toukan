package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Lanes struct {
	lanes  []*tview.List
	active int
	flex   *tview.Flex
	app    *tview.Application
}

func NewLanes(number int, app *tview.Application) *Lanes {
	l := &Lanes{make([]*tview.List, number), 0, tview.NewFlex(), app}
	for i := range l.lanes {
		l.lanes[i] = tview.NewList()
		l.lanes[i].ShowSecondaryText(false).SetBorder(true)
		l.lanes[i].SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyLeft:
				l.DecActive()
				return nil
			case tcell.KeyRight:
				l.IncActive()
				return nil
			}
			switch event.Rune() {
			case 'q':
				app.Stop()
			}
			return event
		})
		for j := 0; j < 10; j++ {
			l.lanes[i].AddItem(fmt.Sprint("Item ", j), "", 0, nil)
		}
		l.flex.AddItem(l.lanes[i], 0, 1, i == 0)
	}

	return l
}

func (l *Lanes) IncActive() {
	l.active++
	l.setActive()
}

func (l *Lanes) DecActive() {
	l.active--
	l.setActive()
}

func (l *Lanes) GetUi() *tview.Flex {
	return l.flex
}

func (l *Lanes) setActive() {
	for l.active < 0 {
		l.active += len(l.lanes)
	}
	l.active %= len(l.lanes)
	l.app.SetFocus(l.lanes[l.active])
}

func main() {
	app := tview.NewApplication()

	lanes := NewLanes(3, app)

	pages := tview.NewPages().
		AddPage("TouKan", lanes.GetUi(), true, true)
	app.SetRoot(pages, true)

	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}
