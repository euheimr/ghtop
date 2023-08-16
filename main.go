package main

// github.com/shirou/gopsutil
// github.com/rivo/tview
// github.com/gdamore/tcell

import (
	"github.com/euheimr/ghtop/ui"
	"github.com/euheimr/ghtop/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"time"
)

var cfg = util.ReadConfig()

var RefreshInterval = cfg.Duration("UpdateInterval") * time.Millisecond
var GroupProcesses = cfg.Bool("GroupProcesses")
var TempScale = cfg.String("TempScale")

var (
	app   *tview.Application
	fMain *tview.Flex
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("[%s] Failed to initialize tcell.NewScreen(): %v",
			util.GetFuncName(), err)
	}
	// fMain is the main box drawn to the screen. It holds all the other boxes
	//	within it.
	fMain = tview.NewFlex().SetDirection(tview.FlexRow)
	fMain.SetBorderStyle(tcell.StyleDefault)
	//fMain.SetBorder(false)

	app = tview.NewApplication().
		SetScreen(screen).
		EnableMouse(true).
		ResizeToFullScreen(fMain)

	sysInfoBox := tview.NewTextView()
	cpuBox := tview.NewBox()
	memBox := tview.NewBox()

	procsTbl := tview.NewTable()
	cpuTempBox := tview.NewBox()
	netBox := tview.NewBox()

	gpuBox := tview.NewBox()
	gpuTempBox := tview.NewBox()

	fRow1 := tview.NewFlex()
	fRow1.AddItem(sysInfoBox, 0, 3, false).
		AddItem(cpuBox, 0, 10, false).
		AddItem(memBox, 0, 6, false)

	fRow2 := tview.NewFlex()
	fRow2.AddItem(procsTbl, 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(cpuTempBox, 0, 1, false).
			AddItem(netBox, 0, 1, false),
			0, 1, false)

	// If there is a GPU, then add `GPU` and `GPU Temp` boxes
	if gpuBox != nil {
		fRow2.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(gpuBox, 0, 1, false).
			AddItem(gpuTempBox, 0, 1, false),
			0, 1, false)
	}
	// Add Row1 and Row2 to flexMain
	fMain.AddItem(fRow1, 0, 1, false)
	fMain.AddItem(fRow2, 0, 2, false)

	// These functions are where all the boxes are drawn via Go Routines
	go ui.UpdateSysInfoBox(app, sysInfoBox, RefreshInterval)
	go ui.UpdateCpuBox(app, cpuBox, RefreshInterval)
	go ui.UpdateMemBox(app, memBox, RefreshInterval)
	go ui.UpdateProcBox(app, procsTbl, RefreshInterval, GroupProcesses)
	go ui.UpdateCpuTempBox(app, cpuTempBox, RefreshInterval)
	// netBox
	// gpuBox
	// gpuTempBox

	if err := app.SetRoot(fMain, true).SetFocus(procsTbl).Run(); err != nil {
		panic(err)
	}
}
