package main

import (
	"fmt"
)

type Content struct {
	items  [][]string
	titles []string
}

func NewContent() *Content {
	ret := &Content{}
	ret.items = make([][]string, 3)
	for i := range ret.items {
		ret.items[i] = make([]string, 10)
		for j := range ret.items[i] {
			ret.items[i][j] = fmt.Sprint("Item ", j)
		}
	}

	ret.titles = []string{"To Do", "Doing", "Done"}

	return ret
}

func (c *Content) GetNumLanes() int {
	return len(c.titles)
}

func (c *Content) GetLaneTitle(idx int) string {
	return c.titles[idx]
}

func (c *Content) GetLaneItems(idx int) []string {
	return c.items[idx]
}

func (c *Content) MoveItem(fromlane, fromidx, tolane, toidx int) {
	item := c.items[fromlane][fromidx]
	// https://github.com/golang/go/wiki/SliceTricks
	c.items[fromlane] = append(c.items[fromlane][:fromidx], c.items[fromlane][fromidx+1:]...)
	c.items[tolane] = append(c.items[tolane][:toidx], append([]string{item}, c.items[tolane][toidx:]...)...)
}

func (c *Content) DelItem(lane, idx int) {
	c.items[lane] = append(c.items[lane][:idx], c.items[lane][idx+1:]...)
}
