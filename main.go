package main

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/rivo/tview"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fname := path.Join(usr.HomeDir, ".toukan.json")

	var content *Content
	f, err := os.Open(fname)
	if err == nil {
		content = NewContentIo(f)
		f.Close()
	}

	if content == nil {
		content = NewContentDefault()
	}

	app := tview.NewApplication()

	lanes := NewLanes(content, app)

	app.SetRoot(lanes.GetUi(), true)

	if err := app.Run(); err != nil {
		log.Fatal("Error running application: %s\n", err)
	}

	f, err = os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	content.Save(f)
	f.Close()
}
