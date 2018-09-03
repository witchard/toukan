package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Lanes struct {
	content  *Content
	lanes    []*tview.List
	active   int
	pages    *tview.Pages
	app      *tview.Application
	inselect bool
}

func NewLanes(content *Content, app *tview.Application) *Lanes {
	l := &Lanes{content, make([]*tview.List, content.GetNumLanes()), 0, tview.NewPages(), app, false}

	flex := tview.NewFlex()
	for i := 0; i < l.content.GetNumLanes(); i++ {
		l.lanes[i] = tview.NewList()
		l.lanes[i].ShowSecondaryText(false).SetBorder(true)
		l.lanes[i].SetTitle(l.content.GetLaneTitle(i))
		l.lanes[i].SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyUp:
				if l.inselect {
					l.up()
					return nil
				}
			case tcell.KeyDown:
				if l.inselect {
					l.down()
					return nil
				}
			case tcell.KeyLeft:
				if l.inselect {
					l.left()
				} else {
					l.decActive()
				}
				return nil
			case tcell.KeyRight:
				if l.inselect {
					l.right()
				} else {
					l.incActive()
				}
				return nil
			}
			switch event.Rune() {
			case 'q':
				l.pages.ShowPage("quit")
			case 'd':
				l.pages.ShowPage("delete")
			case 'a':
				l.pages.ShowPage("add")
			}
			return event
		})
		l.lanes[i].SetSelectedFunc(func(w int, x string, y string, z rune) {
			l.selected()
		})
		l.lanes[i].SetDoneFunc(func() {
			// Cancel select on Done (escape)
			if l.inselect {
				l.selected()
			}
		})
		for _, text := range l.content.GetLaneItems(i) {
			l.lanes[i].AddItem(text, "", 0, nil)
		}
		flex.AddItem(l.lanes[i], 0, 1, i == 0)
	}
	l.pages.AddPage("lanes", flex, true, true)

	quit := tview.NewModal().
		SetText("Do you want to quit the application?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			}
			l.pages.HidePage("quit")
		})
	l.pages.AddPage("quit", quit, false, false)

	delete := tview.NewModal().
		SetText("Are you sure?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				item := l.lanes[l.active].GetCurrentItem()
				l.content.DelItem(l.active, item)
				l.redraw(l.active, item)
			}
			l.pages.HidePage("delete")
		})
	l.pages.AddPage("delete", delete, false, false)

	add := NewModalInput()
	add.SetDoneFunc(func(text string, success bool) {
		if success {
			item := l.lanes[l.active].GetCurrentItem()
			l.content.AddItem(l.active, item, text)
			l.redraw(l.active, item)
		}
		l.pages.HidePage("add")
	})
	l.pages.AddPage("add", add, false, false)

	return l
}

func (l *Lanes) selected() {
	if l.inselect {
		l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorWhite)
	} else {
		l.lanes[l.active].SetSelectedBackgroundColor(tcell.ColorGreen)
	}
	l.inselect = !l.inselect
}

func (l *Lanes) redraw(lane, active int) {
	l.lanes[lane].Clear()
	for _, text := range l.content.GetLaneItems(lane) {
		l.lanes[lane].AddItem(text, "", 0, nil)
	}
	l.lanes[lane].SetCurrentItem(normPos(active, l.lanes[lane].GetItemCount()))
}

func (l *Lanes) up() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newPos := normPos(currentPos-1, l.lanes[l.active].GetItemCount())
	l.content.MoveItem(l.active, currentPos, l.active, newPos)
	l.redraw(l.active, newPos)
}

func (l *Lanes) down() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newPos := normPos(currentPos+1, l.lanes[l.active].GetItemCount())
	l.content.MoveItem(l.active, currentPos, l.active, newPos)
	l.redraw(l.active, newPos)
}

func (l *Lanes) left() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newLane := normPos(l.active-1, len(l.lanes))
	newPos := l.lanes[newLane].GetCurrentItem()
	l.content.MoveItem(l.active, currentPos, newLane, newPos)
	l.redraw(l.active, currentPos)
	l.redraw(newLane, newPos)
	l.selected()
	l.decActive()
	l.selected()
}

func (l *Lanes) right() {
	currentPos := l.lanes[l.active].GetCurrentItem()
	newLane := normPos(l.active+1, len(l.lanes))
	newPos := l.lanes[newLane].GetCurrentItem()
	l.content.MoveItem(l.active, currentPos, newLane, newPos)
	l.redraw(l.active, currentPos)
	l.redraw(newLane, newPos)
	l.selected()
	l.incActive()
	l.selected()

}

func (l *Lanes) decActive() {
	l.active--
	l.setActive()
}

func (l *Lanes) incActive() {
	l.active++
	l.setActive()
}

func normPos(pos, length int) int {
	for pos < 0 {
		pos += length
	}
	pos %= length
	return pos
}

func (l *Lanes) setActive() {
	l.active = normPos(l.active, len(l.lanes))
	l.app.SetFocus(l.lanes[l.active])
}

func (l *Lanes) GetUi() *tview.Pages {
	return l.pages
}
