package ui

import (
	"github.com/rivo/tview"
	"time"
)

var procsTblLabel = "[ Processes ]"

func UpdateProcs(app *tview.Application, procsTbl *tview.Table, update time.Duration) {

	procsTbl.SetBorder(true).SetTitle(procsTblLabel)

	for {
		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})
	}
}
