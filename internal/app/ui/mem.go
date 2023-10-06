package ui

import (
	"github.com/rivo/tview"
	"time"
)

var memTitle = "[ Memory ]"

func UpdateMem(app *tview.Application, mem *tview.Box, update time.Duration) {
	mem.SetBorder(true).SetTitle(memTitle)
	for {
		//get mem data
		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})
	}
}
