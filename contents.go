package main

import (
	"encoding/json"
	"io"
)

type Content struct {
	Titles []string
	Items  [][]string
}

func NewContentIo(r io.Reader) *Content {
	decoder := json.NewDecoder(r)
	c := &Content{}
	if err := decoder.Decode(c); err != nil {
		return nil
	}
	return c
}

func NewContentDefault() *Content {
	ret := &Content{}
	ret.Titles = []string{"To Do", "Doing", "Done"}
	ret.Items = make([][]string, 3)
	return ret
}

func (c *Content) GetNumLanes() int {
	return len(c.Titles)
}

func (c *Content) GetLaneTitle(idx int) string {
	return c.Titles[idx]
}

func (c *Content) GetLaneItems(idx int) []string {
	return c.Items[idx]
}

func (c *Content) MoveItem(fromlane, fromidx, tolane, toidx int) {
	item := c.Items[fromlane][fromidx]
	// https://github.com/golang/go/wiki/SliceTricks
	c.Items[fromlane] = append(c.Items[fromlane][:fromidx], c.Items[fromlane][fromidx+1:]...)
	c.Items[tolane] = append(c.Items[tolane][:toidx], append([]string{item}, c.Items[tolane][toidx:]...)...)
}

func (c *Content) DelItem(lane, idx int) {
	c.Items[lane] = append(c.Items[lane][:idx], c.Items[lane][idx+1:]...)
}

func (c *Content) AddItem(lane, idx int, text string) {
	c.Items[lane] = append(c.Items[lane][:idx], append([]string{text}, c.Items[lane][idx:]...)...)
}

func (c *Content) Save(w io.Writer) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(c)
}
