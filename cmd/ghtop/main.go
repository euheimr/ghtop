package main

import (
	"github.com/euheimr/ghtop/internal/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rivo/tview"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type AppLayout struct {
	info     *tview.TextView
	cpu      *tview.Box
	cpuTemp  *tview.Box
	mem      *tview.Box
	procsTbl *tview.Table
	net      *tview.Box
	gpu      *tview.Box
	gpuTemp  *tview.Box
}

type ConfigVars struct {
	Debug          bool
	UpdateInterval time.Duration
	Celsius        bool
	EnableNvidia   bool
	EnableTUI      bool
}

const (
	CONFIG_FILENAME             = "cfg.toml"
	CONFIG_UPDATE_DELAY_SECONDS = 3
)

var (
	app          *tview.Application
	layout       AppLayout
	views        *[]AppLayout
	selectedView int
)

var cfg = &ConfigVars{
	Debug:          false,
	UpdateInterval: 100 * time.Millisecond,
	Celsius:        true,
	EnableNvidia:   false,
	EnableTUI:      true,
}

func setupLayout(app *tview.Application) {
	views = &[]AppLayout{
		0: {
			// row 1
			info:    tview.NewTextView(),
			cpu:     tview.NewBox(),
			cpuTemp: tview.NewBox(),
			// row 2
			mem:      tview.NewBox(),
			procsTbl: tview.NewTable(),
			net:      tview.NewBox(),
		},
		//1: {
		//	info: tview.NewTextView(),
		//},
	}
	layout = (*views)[0]
	slog.Debug("views count = " + strconv.FormatInt(int64(len(*views)), 10))

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
		AddItem(layout.procsTbl, 0, 2, false).
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

	fMain := tview.NewFlex()
	fMain.
		AddItem(flexRow1, 0, 22, false).
		AddItem(flexRow2, 0, 40, false)
	// this sets the first "Main" layout view to always be rows
	fMain.SetDirection(tview.FlexRow)

	// finally set the root object
	app.SetRoot(fMain, true).EnableMouse(true)
}

func updateConfigVars(k *koanf.Koanf, f *file.File) {
	if err := k.Load(f, toml.Parser()); err != nil {
		slog.Error("Could not load config file! " + err.Error())
	}
	cfg = &ConfigVars{
		Debug:          k.Bool("Debug"),
		UpdateInterval: k.Duration("UpdateInterval") * time.Millisecond,
		Celsius:        k.Bool("Celsius"),
		// TODO: detect AMD / nvidia gpus automatically and override??
		EnableNvidia: k.Bool("EnableNvidia"),
		EnableTUI:    k.Bool("EnableTUI"),
	}
	slog.Info("Updated configuration variables")
}

func startApp(app *tview.Application) {
	// we must first setup the UI layout before starting the goroutines below
	setupLayout(app)

	// queue the draw updates with goroutines
	go ui.UpdateCpu(app, layout.cpu, cfg.UpdateInterval)
	go ui.UpdateCpuTemp(app, layout.cpuTemp, cfg.UpdateInterval)
	go ui.UpdateMem(app, layout.mem, cfg.UpdateInterval)
	go ui.UpdateNet(app, layout.net, cfg.UpdateInterval)
	go ui.UpdateProcs(app, layout.procsTbl, cfg.UpdateInterval)
	if cfg.EnableNvidia {
		go ui.UpdateGpu(app, layout.gpu, cfg.UpdateInterval)
		go ui.UpdateGpuTemp(app, layout.gpuTemp, cfg.UpdateInterval)
	}

	// We set the keybinds here (Quit app, force reload, change view, etc ...)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
		//case tcell.KeyEsc:
		//	app.Stop()
		default:
			return event
		}
		return event
	})

	// Finally, run the app!
	if err := app.Run(); err != nil {
		slog.Error("Application error! " + err.Error())
		os.Exit(1)
	}
}

func main() {
	app = tview.NewApplication()

	slog.Debug("Loading cfg default values ...")
	k := koanf.New(".")
	f := file.Provider(CONFIG_FILENAME)

	// If the config file exists, update `cfg` using updateConfigVars()
	if _, err := os.Stat(CONFIG_FILENAME); err == nil {
		// load config values from file and start the app
		updateConfigVars(k, f)
		startApp(app)

		// also watch for any file changes and restart the app as needed
		f.Watch(func(event interface{}, err error) {
			if err != nil {
				slog.Error("Cannot watch for config file changes! " + err.Error())
			}
			time.Sleep(CONFIG_UPDATE_DELAY_SECONDS * time.Second)
			app.Suspend(func() {
				// load the new values and restart
				updateConfigVars(k, f)
				startApp(app)
			})
		})
	} else {
		// this code is run if a config file does not exist, effectively using default values
		if cfg.EnableTUI {
			// If the text UI is enabled, run the app. Otherwise, don't start it.
			//	This is mostly for debugging. Eventually I'll log to file... but not today
			startApp(app)
		} else {
			slog.Info("Did not start app - EnableTUI is " +
				strconv.FormatBool(cfg.EnableTUI) + " !")
		}
	}
}
