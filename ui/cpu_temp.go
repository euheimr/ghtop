package ui

import (
	"github.com/rivo/tview"
	"time"
)

var cpuTempBoxLabel = "[ CPU Temp ]"

func UpdateCpuTempBox(app *tview.Application, cpuTempBox *tview.Box,
	refresh time.Duration) {

	for {
		// TODO: get cpu temp data

		time.Sleep(refresh)
		app.QueueUpdateDraw(func() {
			// TODO: draw the braille graph

			// TODO: draw the temp data text below the graph

		})

		// draw border and title last
		cpuTempBox.SetBorder(true).SetTitle(cpuTempBoxLabel)
	}
}
