package main

import (
	"fmt"

	"github.com/gizak/termui"
)

// DisplayList represents a list of items within the UI.
type DisplayList struct {
	List         *termui.List
	items        []string
	name         string
	selected     int
	selectedOrig string
}

// NewDisplayList creates a new DisplayList structure with no items.
func NewDisplayList(name string, items []string) *DisplayList {
	l := DisplayList{termui.NewList(), items, name, -1, ""}
	l.List.Items = l.items
	l.List.BorderLabel = l.name
	l.Resize()
	return &l
}

// SetItems is used to set what items are in the list.
func (l *DisplayList) SetItems(items []string) {
	l.selected = -1 // Reset selected position
	l.items = items
	l.List.Items = l.items
}

// Resize is used to update the item on a terminal resize.
func (l *DisplayList) Resize() {
	l.List.Height = termui.TermHeight()
}

// Select is used to highlight an item. Negative values indicate nothing selected.
func (l *DisplayList) Select(pos int) {
	// Reset selected item
	if l.selected >= 0 {
		l.items[l.selected] = l.selectedOrig
	}
	l.selected = -1

	// Compute new pos and update selected
	if pos >= 0 {
		l.selected = pos % len(l.items)
		l.selectedOrig = l.items[l.selected]
		l.items[l.selected] = fmt.Sprintf("[%s](fg-black,bg-white)", l.selectedOrig)
	}
}

// IncSelected is used to increment the selected item in the list.
func (l *DisplayList) IncSelected() {
	l.Select(l.selected + 1)
}

// DecSelected is used to decrement the selected item in the list.
func (l *DisplayList) DecSelected() {
	l.Select(len(l.items) + l.selected - 1) // Cope with zero wrap by adding length - this is corrected in Select function
}

// GetSelected is used to retrieve the currently selected item.
func (l *DisplayList) GetSelected() int {
	return l.selected
}

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	list1 := NewDisplayList("To Do", []string{"foo", "bar", "baz"})
	list2 := NewDisplayList("Doing", []string{"spam", "egg", "chips"})
	list3 := NewDisplayList("Done", []string{"gopher"})
	list1.Select(0)

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(4, 0, list1.List),
			termui.NewCol(4, 0, list2.List),
			termui.NewCol(4, 0, list3.List),
		),
	)
	termui.Body.Align()
	termui.Render(termui.Body)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<down>", func(e termui.Event) {
		list1.IncSelected()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/kbd/<up>", func(e termui.Event) {
		list1.DecSelected()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		list1.Resize()
		list2.Resize()
		list3.Resize()
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Loop()
}
