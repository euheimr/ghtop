package ui

import (
	"github.com/rivo/tview"
	"time"
)

var gpuBoxLabel = "[ GPU ]"
var gpuTempBoxLabel = "[ GPU Temp ]"

func UpdateGpuBox(app *tview.Application, gpuBox *tview.Box,
	update time.Duration) {

	gpuBox.SetBorder(true).SetTitle(gpuBoxLabel)

	for {
		// todo: get gpu data

		time.Sleep(update)
		app.QueueUpdateDraw(func() {
			//todo: draw braille graph

			// todo: draw gpu text below the graph

		})
	}
}

func UpdateGpuTempBox(app *tview.Application, gpuTempBox *tview.Box,
	update time.Duration) {

	gpuTempBox.SetBorder(true).SetTitle(gpuTempBoxLabel)

	for {
		// todo: get gpu data

		time.Sleep(update)
		app.QueueUpdateDraw(func() {
			//todo: draw braille graph

			// todo: draw gpu text below the graph

		})
	}
}
