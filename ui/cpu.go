package ui

import (
	"github.com/euheimr/ghtop/devices"
	"github.com/rivo/tview"
	"time"
)

var cpuModelName = devices.CpuModelName()
var cpuBoxLabel = "[ " + cpuModelName + " ]"
var cpuTempBoxLabel = "[ CPU Temp ]"

func UpdateCpuBox(app *tview.Application, cpuBox *tview.Box,
	update time.Duration) {

	cpuBox.SetBorder(true).SetTitle(cpuBoxLabel)

	for {
		// TODO: get cpu data

		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})
	}
}

func UpdateCpuTempBox(app *tview.Application, cpuTempBox *tview.Box,
	update time.Duration) {

	cpuTempBox.SetBorder(true).SetTitle(cpuTempBoxLabel)

	for {
		// TODO: get cpu temp data

		time.Sleep(update)
		app.QueueUpdateDraw(func() {
			// TODO: draw the braille graph

			// TODO: draw the temp data text below the graph

		})

	}
}
