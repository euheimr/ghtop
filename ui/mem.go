package ui

import (
	"github.com/rivo/tview"
	"time"
)

var memBoxLabel = "[ Memory ]"

func UpdateMemBox(app *tview.Application, memBox *tview.Box,
	update time.Duration) {

	memBox.SetBorder(true).SetTitle(memBoxLabel)

	for {
		// TODO: get memory data and return it with the box

		time.Sleep(update)
		app.QueueUpdateDraw(func() {
			// TODO: draw the braille graph

			// TODO: draw the memory data text below the graph
		})

	}

}
