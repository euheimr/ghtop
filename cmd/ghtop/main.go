package main

import (
	"github.com/euheimr/ghtop/internal"
	"github.com/euheimr/ghtop/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"os"
	"strconv"
)

type AppLayout struct {
	info    *tview.TextView
	cpu     *tview.Box
	cpuTemp *tview.Box
	mem     *tview.Box
	procTbl *tview.Table
	net     *tview.Box
	gpu     *tview.Box
	gpuTemp *tview.Box
}

var (
	app          *tview.Application
	cfg          *internal.ConfigVars
	layout       AppLayout
	views        *[]AppLayout
	selectedView int
)

const ENABLE_APP = true

func setupKeybinds(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
		case tcell.KeyEsc:
			app.Stop()
		}
		return event
	})
}

func setupLayout(app *tview.Application) {
	views = &[]AppLayout{
		0: {
			// row 1
			info:    tview.NewTextView(),
			cpu:     tview.NewBox(),
			cpuTemp: tview.NewBox(),
			// row 2
			mem:     tview.NewBox(),
			procTbl: tview.NewTable(),
			net:     tview.NewBox(),
		},
		1: {
			info: tview.NewTextView(),
		},
	}
	layout = (*views)[0]

	slog.Info("views count = " + strconv.FormatInt(int64(len(*views)), 10))

	// build row 1
	flexRow1 := tview.NewFlex()
	if cfg.Debug {
		flexRow1.
			AddItem(layout.info, 0, 2, false).
			AddItem(layout.cpu, 0, 7, false).
			AddItem(layout.mem, 0, 3, false)
	} else {
		flexRow1.
			AddItem(layout.cpu, 0, 3, false).
			AddItem(layout.mem, 0, 1, false)
	}

	//build row 2
	flexRow2 := tview.NewFlex()
	flexRow2.
		// row 2 column 1
		AddItem(layout.procTbl, 0, 2, false).
		// row 2 column 2
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(layout.cpuTemp, 0, 2, false).
			AddItem(layout.net, 0, 2, false),
			0, 1, false)

	// if theres a GPU then add `GPU` and `GPUTemp` boxes
	if cfg.EnableNvidia {
		layout.gpu = tview.NewBox()
		layout.gpuTemp = tview.NewBox()

		flexRow2.
			// row 2 column 3
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(layout.gpu, 0, 4, false).
				AddItem(layout.gpuTemp, 0, 4, false),
				0, 1, false)
	}
	// todo: add row1 and row2??
	fMain := tview.NewFlex()
	fMain.
		AddItem(flexRow1, 0, 22, false).
		AddItem(flexRow2, 0, 40, false)
	// this sets the first "Main" layout view to always be rows
	fMain.SetDirection(tview.FlexRow)

	// finally set the root object
	app.SetRoot(fMain, true).EnableMouse(true)
}

func startDraw(app *tview.Application) {
	setupLayout(app)

	go ui.UpdateCpu(app, layout.cpu, cfg.UpdateInterval)
	go ui.UpdateCpuTemp(app, layout.cpuTemp, cfg.UpdateInterval)
	go ui.UpdateMem(app, layout.mem, cfg.UpdateInterval)
	go ui.UpdateNet(app, layout.net, cfg.UpdateInterval)
	go ui.UpdateProc(app, layout.procTbl, cfg.UpdateInterval)

}

func main() {
	app = tview.NewApplication()
	cfg = internal.Cfg

	setupKeybinds(app)
	startDraw(app)

	if ENABLE_APP {
		if err := app.Run(); err != nil {
			slog.Error("Application error! ", err.Error())
			os.Exit(1)
		}
	}
}
