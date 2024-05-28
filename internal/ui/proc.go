package ui

import (
	"github.com/rivo/tview"
	"time"
)

var procTblLabel = "[ Processes ]"

func UpdateProc(app *tview.Application, procTbl *tview.Table, update time.Duration) {

	procTbl.SetBorder(true).SetTitle(procTblLabel)

	for {
		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})
	}
}
