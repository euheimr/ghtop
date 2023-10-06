package ui

import (
	"github.com/rivo/tview"
	"time"
)

var (
	gpuLabel     = "[ GPU ]"
	gpuTempLabel = "[ GPU Temp ]"
)

func UpdateGpu(app *tview.Application, gpu *tview.Box,
	update time.Duration) {

	gpu.SetBorder(true).SetTitle(gpuLabel)

	for {

		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})

	}
}

func UpdateGpuTemp(app *tview.Application, gpuTemp *tview.Box,
	update time.Duration) {

	gpuTemp.SetBorder(true).SetTitle(gpuTempLabel)

	for {

		time.Sleep(update)
		app.QueueUpdateDraw(func() {

		})

	}
}
