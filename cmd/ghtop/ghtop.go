package main

import (
	"github.com/euheimr/ghtop/internal"
	"github.com/euheimr/ghtop/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log/slog"
	"os"
	"time"
)

type RowButtons struct {
	spacer  *tview.TextView
	btn0    *tview.Button
	btn1    *tview.Button
	btn2    *tview.Button
	opt     *tview.Button
	dbg     *tview.Button
	btnExit *tview.Button
}

type ModulesMain struct {
	btns     *RowButtons
	sysInfo  *tview.TextView
	cpu      *tview.Box
	mem      *tview.Box
	procsTbl *tview.Table
	cpuTemp  *tview.Box
	net      *tview.Box
	gpu      *tview.Box
	gpuTemp  *tview.Box
}

type State struct {
	app          *tview.Application
	cfg          *internal.ConfigVars
	log          *slog.Logger
	selectedView int
	views        []*tview.Flex
	//defaultView  int
}

var AppState *State
var Modules *ModulesMain

func init() {
	// Setup the initial state for ghtop
	AppState = &State{
		app: tview.NewApplication(),
		cfg: &internal.Config,
		views: []*tview.Flex{
			//	>flexMain holds all the other boxes within it.<
			0: tview.NewFlex(),
			1: tview.NewFlex(),
			//2: tview.NewFlex(),
		},
		//defaultView:  0,
	}
	AppState.selectedView = AppState.cfg.SelectedViewOverride

	AppState.setupLogging()

	// this sets the first "Main layout View" to always be rows
	// TODO: have some way of iterating over views[]?
	AppState.views[0].SetDirection(tview.FlexRow)
	AppState.views[1].SetDirection(tview.FlexRow)

}

func (s *State) setupLogging() {
	/// Start log handling	////////////////////////////////////////////////////
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		//AddSource:   true,
		//ReplaceAttr: nil,
	}
	if !s.cfg.Debug {
		opts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}

	// log to file
	f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("Failed to open log file")
	}
	defer f.Close()

	logger := slog.New(slog.NewTextHandler(f, opts)) //os.Stdout, opts))
	slog.SetDefault(logger)

	s.log = slog.Default()
	s.log.Debug(internal.GetFuncName() + "Initialized logging")
}

func (s *State) setupKeybinds() {
	s.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			s.app.Stop()
		case tcell.KeyEsc:
			s.app.Stop()
			//case tcell.KeyCtrlQ:
			//	a.Stop()
		}
		return event
	})
}

func (s *State) addButtonRowToViews() {
	Modules = &ModulesMain{
		btns: &RowButtons{
			btn0:    tview.NewButton("btn0"),
			btn1:    tview.NewButton("btn1"),
			btn2:    tview.NewButton("btn2"),
			spacer:  tview.NewTextView(),
			dbg:     tview.NewButton("DBG"),
			opt:     tview.NewButton("OPT"),
			btnExit: tview.NewButton("EXIT"),
		},
	}

	// tabs buttons and spacer
	Modules.btns.btnExit.SetBackgroundColor(tcell.ColorRed)
	fRowButtons := tview.NewFlex().SetDirection(tview.FlexColumn)
	fRowButtons.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(Modules.btns.btn0, 0, 1, false).
		AddItem(Modules.btns.btn1, 0, 1, false).
		AddItem(Modules.btns.btn2, 0, 1, false).
		AddItem(Modules.btns.spacer, 0, 8, false).
		AddItem(Modules.btns.dbg, 0, 1, false).
		AddItem(Modules.btns.opt, 0, 1, false).
		AddItem(Modules.btns.btnExit, 0, 1, false),
		0, 1, false)

	// TODO: find some way of iterating views in the future via index?
	s.views[0].AddItem(fRowButtons, 0, 2, false)
	s.views[1].AddItem(fRowButtons, 0, 2, false)
	s.log.Debug(internal.GetFuncName() + "Init addButtonRowToViews()")
}

func (s *State) setupLayoutMain() {
	Modules = &ModulesMain{
		// row1
		sysInfo: tview.NewTextView(),
		cpu:     tview.NewBox(),
		mem:     tview.NewBox(),
		// row2
		procsTbl: tview.NewTable(),
		cpuTemp:  tview.NewBox(),
		net:      tview.NewBox(),
	}

	fRow1 := tview.NewFlex()
	if s.cfg.Debug {
		fRow1.
			AddItem(Modules.sysInfo, 0, 2, false).
			AddItem(Modules.cpu, 0, 7, false).
			AddItem(Modules.mem, 0, 3, false)
	} else {
		fRow1.
			AddItem(Modules.cpu, 0, 3, false).
			AddItem(Modules.mem, 0, 1, false)
	}

	fRow2 := tview.NewFlex()
	fRow2.AddItem(Modules.procsTbl, 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(Modules.cpuTemp, 0, 2, false).
			AddItem(Modules.net, 0, 2, false),
			0, 1, false)

	// If there is a GPU, then add `GPU` and `GPU Temp` boxes
	if s.cfg.EnableNvidia {
		Modules.gpu = tview.NewBox()
		Modules.gpuTemp = tview.NewBox()

		fRow2.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(Modules.gpu, 0, 4, false).
			AddItem(Modules.gpuTemp, 0, 4, false),
			0, 1, false)
	}

	s.views[0].AddItem(fRow1, 0, 22, false)
	s.views[0].AddItem(fRow2, 0, 40, false)
	slog.Debug("" + internal.GetFuncName() + "Init setupLayoutMain()")
}

func (s *State) startDraw() {
	switch s.selectedView {
	case 0:
		// todo: start the goroutines
		// These functions are where all the boxes are drawn via Go Routines
		if Modules.sysInfo != nil {
			go ui.UpdateSysInfo(s.app, Modules.sysInfo, s.cfg.UpdateInterval)
		}
		if Modules.cpu != nil {
			go ui.UpdateCpu(s.app, Modules.cpu, s.cfg.UpdateInterval)
		}
		if Modules.mem != nil {
			go ui.UpdateMem(s.app, Modules.mem, s.cfg.UpdateInterval)
		}
		if Modules.procsTbl != nil {
			go ui.UpdateProcs(s.app, Modules.procsTbl, s.cfg.UpdateInterval)
		}

		if Modules.cpuTemp != nil {
			go ui.UpdateCpuTemp(s.app, Modules.cpuTemp, s.cfg.UpdateInterval)
		}

		if Modules.net != nil {
			go ui.UpdateNet(s.app, Modules.net, s.cfg.UpdateInterval)
		}

		if (Modules.gpu != nil || Modules.gpuTemp != nil) && s.cfg.EnableNvidia {
			go ui.UpdateGpu(s.app, Modules.gpu, s.cfg.UpdateInterval)
			go ui.UpdateGpuTemp(s.app, Modules.gpuTemp, s.cfg.UpdateInterval)
		}

		//if Modules.btns != nil {
		//go UpdateButtons(s.app, Modules.btns, s.cfg.UpdateInterval)
		//}
	}
}

func main() {
	AppState.setupKeybinds()
	if AppState.cfg.EnableUIButtons {
		AppState.addButtonRowToViews()
	}
	AppState.setupLayoutMain()
	AppState.startDraw()

	/// !TODO: testing some alternate views
	var (
		r1Proportion = 45
		m            = tview.NewFlex()
		r1           = tview.NewFlex()
		r2           = tview.NewFlex()
	)
	r1.SetBorder(true).SetBackgroundColor(tcell.ColorGray)
	r2.SetBorder(true).SetBackgroundColor(tcell.ColorOrange)
	m.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(r1, 0, 2, false).
		AddItem(r2, 0, 2, false),
		0, 1, false)
	AppState.views[1].AddItem(m, 0, r1Proportion, false)
	//AppState.views[1].AddItem(r2, 0, r2Prop, false)
	/// END !TODO: testing some alternate views

	switch AppState.selectedView {
	default:
		AppState.app.SetRoot(AppState.views[0], true).
			//SetFocus(Modules.procsTbl).
			EnableMouse(true)
		AppState.log.Debug("Active view set to: 0")
	case 0:
		AppState.app.SetRoot(AppState.views[0], true).
			//SetFocus(Modules.procsTbl).
			EnableMouse(true)
		AppState.log.Debug("Active view set to: 0")
	case 1:
		AppState.app.SetRoot(AppState.views[1], true).
			//SetFocus(Modules.procsTbl).
			EnableMouse(true)
		AppState.log.Debug("Active view set to: 1")
	}

	if err := AppState.app.Run(); err != nil {
		//panic(err)
		AppState.log.Error("Application error: ", err)
		os.Exit(1)
	}
}
