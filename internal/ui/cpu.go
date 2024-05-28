package ui

import (
	"github.com/rivo/tview"
	"time"
)

var cpuLabel = "[ " + "" + " ]"
var cpuTempLabel = "[ CPU Temp ]"

func UpdateCpu(app *tview.Application, module *tview.Box, update time.Duration) {

	module.SetBorder(true).SetTitle(cpuLabel)
	for {
		// TODO: get cpu data before the sleep
		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})
	}
}

func UpdateCpuTemp(app *tview.Application, module *tview.Box, update time.Duration) {

	module.SetBorder(true).SetTitle(cpuTempLabel)
	for {
		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})
	}
}
