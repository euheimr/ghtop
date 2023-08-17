package ui

import (
	"github.com/rivo/tview"
	"time"
)

var netBoxLabel = "[ Network ]"

func UpdateNetBox(app *tview.Application, netBox *tview.Box,
	update time.Duration) {

	netBox.SetBorder(true).SetTitle(netBoxLabel)

	for {

		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})

	}
}
