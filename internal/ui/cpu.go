package ui

import (
	"github.com/euheimr/ghtop/internal/devices"
	"github.com/rivo/tview"
	"time"
)

type CpuData struct {
}

var cpuLabel = "[ " + devices.CpuInfo[0].ModelName + " ]"
var cpuTempLabel = "[ CPU Temp ]"

func init() {
	// todo: Get initial CPU metrics/data

	// todo: Get initial CPU Temp data

}

func UpdateCpu(app *tview.Application, cpu *tview.Box, update time.Duration) {
	cpu.SetBorder(true).SetTitle(cpuLabel)
	for {
		// get cpu data
		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})
	}
}

func UpdateCpuTemp(app *tview.Application, cpuTemp *tview.Box, update time.Duration) {

	cpuTemp.SetBorder(true).SetTitle(cpuTempLabel)
	for {
		// get cpu data
		time.Sleep(update)
		app.QueueUpdateDraw(func() {
			// TODO: draw the braille graph

			// TODO: draw the temp data text below the graph
		})
	}
}
