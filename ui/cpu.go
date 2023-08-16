package ui

import (
	"github.com/euheimr/ghtop/devices"
	"github.com/rivo/tview"
	"time"
)

var cpuModelName = devices.CpuModelName()
var cpuBoxLabel = "[ " + cpuModelName + " ]"

func UpdateCpuBox(app *tview.Application, cpuBox *tview.Box,
	refresh time.Duration) {

	cpuBox.SetBorder(true).SetTitle(cpuBoxLabel)

	for {
		// TODO: get cpu data

		time.Sleep(refresh)
		app.QueueUpdateDraw(func() {

		})
	}
}
