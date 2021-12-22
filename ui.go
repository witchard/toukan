package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Lanes struct {
	content  *Content
	lanes    []*tview.List
	active   int
	pages    *tview.Pages
	app      *tview.Application
	inselect bool
	add      *ModalInput
	edit     *ModalInput
}

func NewLanes(content *Content, app *tview.Application) *Lanes {
	l := &Lanes{content, make([]*tview.List, content.GetNumLanes()), 0, tview.NewPages(), app, false, NewModalInput(), NewModalInput()}

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
				l.add.SetValue("")
				l.pages.ShowPage("add")
			case 'e':
				if item := l.currentItem(); item != nil {
					l.edit.SetValue(item.Title)
					l.pages.ShowPage("edit")
				}
			case 'n':
				app.Suspend(l.editNote)
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
		for _, item := range l.content.GetLaneItems(i) {
			l.lanes[i].AddItem(item.Title, "", 0, nil)
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
			l.setActive()
		})
	l.pages.AddPage("delete", delete, false, false)

	l.add.SetDoneFunc(func(text string, success bool) {
		if success {
			item := l.lanes[l.active].GetCurrentItem()
			l.content.AddItem(l.active, item, text)
			l.redraw(l.active, item)
		}
		l.pages.HidePage("add")
		l.setActive()
	})
	l.pages.AddPage("add", l.add, false, false)

	l.edit.SetDoneFunc(func(text string, success bool) {
		if success {
			item := l.lanes[l.active].GetCurrentItem()
			itemVal := l.currentItem()
			itemVal.Title = text
			l.redraw(l.active, item)
		}
		l.pages.HidePage("edit")
		l.setActive()
	})
	l.pages.AddPage("edit", l.edit, false, false)

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
	for _, item := range l.content.GetLaneItems(lane) {
		l.lanes[lane].AddItem(item.Title, "", 0, nil)
	}
	num := l.lanes[lane].GetItemCount()
	if num > 0 {
		l.lanes[lane].SetCurrentItem(normPos(active, num))
	}
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

func (l *Lanes) currentItem() *Item {
	pos := l.lanes[l.active].GetCurrentItem()
	content := l.content.GetLaneItems(l.active)
	if pos < 0 || pos >= len(content) {
		return nil
	}
	return &content[pos]
}

func (l *Lanes) editNote() {
	item := l.currentItem()
	if item != nil {
		tmp, err := ioutil.TempFile("", "toukan")
		if err == nil {
			name := tmp.Name()
			defer os.Remove(name)
			tmp.Write([]byte(item.Note))
			tmp.Close()
			cmd := exec.Command("vim", name)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err == nil {
				note_raw, err := ioutil.ReadFile(name)
				if err == nil {
					item.Note = string(note_raw)
				}
			}
		}
	}
}

func (l *Lanes) GetUi() *tview.Pages {
	return l.pages
}
