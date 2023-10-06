package main

import (
	"fmt"
	"github.com/euheimr/ghtop/internal/app/common"
	"github.com/euheimr/ghtop/internal/app/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

type RowButtons struct {
	btn0    *tview.Button
	btn1    *tview.Button
	btn2    *tview.Button
	btn3    *tview.Button
	spacer  *tview.Box
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
	cfg          *common.ConfigVars
	selectedView int
	views        []*tview.Flex
	//defaultView  int
}

var AppState *State
var Log = *common.Logger

func init() {
	// Setup the initial state for ghtop
	AppState = &State{
		app:          tview.NewApplication(),
		cfg:          &common.Config,
		selectedView: 0,
		views: []*tview.Flex{
			//	>flexMain holds all the other boxes within it.<
			0: tview.NewFlex(),
			1: tview.NewFlex(),
			2: tview.NewFlex(),
		},
		//defaultView:  0,
	}
	// this sets the first "Main layout View" to always be rows
	// TODO: have some way of iterating over views[]?
	AppState.views[0].SetDirection(tview.FlexRow)
	AppState.views[1].SetDirection(tview.FlexRow)

	Log.Info("init main")

	//AppState.app.ResizeToFullScreen(AppState.flexMain)

	//for view := range AppState.views { }
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

func (s *State) addButtonRowsToViews() {
	modules := ModulesMain{
		// row0
		btns: &RowButtons{
			btn0:    tview.NewButton("btn0"),
			btn1:    tview.NewButton("btn1"),
			btn2:    tview.NewButton("btn2"),
			btn3:    tview.NewButton("OPT"),
			spacer:  tview.NewBox().SetBackgroundColor(tcell.ColorRed),
			btnExit: tview.NewButton("EXIT"),
		},
	}

	// tabs buttons and spacer
	btns := *modules.btns
	btns.btn0.SetSelectedFunc(func() {
		s.selectedView = 1
		btns.btn0.SetBackgroundColor(tcell.ColorGray)
	})

	btns.btn1.SetSelectedFunc(func() {
		btns.btn1.SetBackgroundColor(tcell.ColorGray).SetTitleColor(tcell.ColorYellow)
	})

	btns.btn2.SetSelectedFunc(func() {
		btns.btn2.SetBackgroundColor(tcell.ColorGray)
	})

	modules.btns.btn3.SetSelectedFunc(func() {
		btns.btn3.SetBackgroundColor(tcell.ColorGray)
	})

	btns.btnExit.SetSelectedFunc(func() {
		s.app.Stop()
	})
	btns.btnExit.SetBackgroundColor(tcell.ColorRed)
	fRowButtons := tview.NewFlex().SetDirection(tview.FlexColumn)
	fRowButtons.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(btns.btn0, 0, 1, false).
		AddItem(btns.btn1, 0, 1, false).
		AddItem(btns.btn2, 0, 1, false).
		AddItem(btns.spacer, 0, 8, false).
		AddItem(btns.btn3, 0, 1, false).
		AddItem(btns.btnExit, 0, 1, false),
		0, 1, false)

	// TODO: find some way of iterating views in the future via index?
	s.views[0].AddItem(fRowButtons, 0, 2, false)
	s.views[1].AddItem(fRowButtons, 0, 2, false)
}

func (s *State) setupLayoutMain() {
	modules := ModulesMain{
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
	fRow1.AddItem(modules.sysInfo, 0, 4, false).
		AddItem(modules.cpu, 0, 14, false).
		AddItem(modules.mem, 0, 6, false)

	fRow2 := tview.NewFlex()
	fRow2.AddItem(modules.procsTbl, 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(modules.cpuTemp, 0, 4, false).
			AddItem(modules.net, 0, 4, false),
			0, 1, false)

	// If there is a GPU, then add `GPU` and `GPU Temp` boxes
	if s.cfg.EnableNvidia {
		modules.gpu = tview.NewBox()
		modules.gpuTemp = tview.NewBox()

		fRow2.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(modules.gpu, 0, 4, false).
			AddItem(modules.gpuTemp, 0, 4, false),
			0, 1, false)
	}

	s.views[0].AddItem(fRow1, 0, 22, false)
	s.views[0].AddItem(fRow2, 0, 40, false)

	// todo: start the goroutines
	// These functions are where all the boxes are drawn via Go Routines
	go ui.UpdateSysInfo(s.app, modules.sysInfo, s.cfg.UpdateInterval)
	go ui.UpdateCpu(s.app, modules.cpu, s.cfg.UpdateInterval)
	go ui.UpdateMem(s.app, modules.mem, s.cfg.UpdateInterval)
	go ui.UpdateProcs(s.app, modules.procsTbl, s.cfg.UpdateInterval)
	go ui.UpdateCpuTemp(s.app, modules.cpuTemp, s.cfg.UpdateInterval)
	go ui.UpdateNet(s.app, modules.net, s.cfg.UpdateInterval)
	if s.cfg.EnableNvidia {
		go ui.UpdateGpu(s.app, modules.gpu, s.cfg.UpdateInterval)
		go ui.UpdateGpuTemp(s.app, modules.gpuTemp, s.cfg.UpdateInterval)
	}

}

func (s *State) initApp(view int) {
	s.app.SetRoot(s.views[view], true).
		//SetFocus(AppState.views[0].modules.procsTbl).
		EnableMouse(true)

	err := s.app.Run()
	if err != nil {
		//panic(err)
		fmt.Print("Application error: ", err)
		os.Exit(1)
	}
}

func main() {

	AppState.setupKeybinds()
	AppState.addButtonRowsToViews()

	AppState.setupLayoutMain()

	// todo: testing some alternate views
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
	//AppState.views[0].AddItem(r2, 0, r2Prop, false)

	AppState.selectedView = 0
	switch AppState.selectedView {
	case 0:
		AppState.initApp(0)
	case 1:
		AppState.initApp(1)
	}

}
