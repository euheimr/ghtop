package ui

import (
	"github.com/rivo/tview"
	"time"
)

var memBoxLabel = "[ Memory ]"

func UpdateMemBox(app *tview.Application, memBox *tview.Box,
	refresh time.Duration) {

	for {
		// TODO: get memory data and return it with the box

		time.Sleep(refresh)
		app.QueueUpdateDraw(func() {
			// TODO: draw the braille graph

			// TODO: draw the memory data text below the graph
		})

		// draw border LAST!
		memBox.SetBorder(true).SetTitle(memBoxLabel)
	}

}
