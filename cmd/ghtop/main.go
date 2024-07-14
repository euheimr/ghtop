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
	GroupProcesses bool
	EnableGPU      bool
	EnableUI       bool
}

const CONFIG_FILENAME = "config.toml"

var (
	app          *tview.Application
	layout       AppLayout
	views        *[]AppLayout
	selectedView int
)

var Cfg = &ConfigVars{
	Debug:          false,
	UpdateInterval: 100 * time.Millisecond,
	Celsius:        true,
	GroupProcesses: true,
	EnableGPU:      false,
	EnableUI:       false,
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

	if Cfg.Debug {
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
	if Cfg.EnableGPU {
		layout.gpu = tview.NewBox()
		layout.gpuTemp = tview.NewBox()

		flexRow2.
			// row 2 column 3
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(layout.gpu, 0, 4, false).
				AddItem(layout.gpuTemp, 0, 4, false),
				0, 1, false)
	}

	row3 := tview.NewTextView()
	row3.SetText(" <F1> Test   <F2> Test 1   <F3> Test 2   <F4> Test 3")

	fMain := tview.NewFlex()
	fMain.AddItem(flexRow1, 0, 22, false)
	fMain.AddItem(flexRow2, 0, 40, false)
	fMain.AddItem(row3, 0, 1, false)
	// this sets the first "Main" layout view to always be rows
	fMain.SetDirection(tview.FlexRow)

	// finally set the root object
	app.SetRoot(fMain, true).EnableMouse(true)
}

func startApp(app *tview.Application) {
	// we must first setup the UI layout before starting the goroutines below
	setupLayout(app)

	// queue the draw updates with goroutines
	go ui.UpdateCpu(app, layout.cpu, Cfg.UpdateInterval)
	go ui.UpdateCpuTemp(app, layout.cpuTemp, Cfg.UpdateInterval)
	go ui.UpdateMem(app, layout.mem, Cfg.UpdateInterval)
	go ui.UpdateNet(app, layout.net, Cfg.UpdateInterval)
	go ui.UpdateProcs(app, layout.procsTbl, Cfg.GroupProcesses, Cfg.UpdateInterval)
	if Cfg.EnableGPU {
		go ui.UpdateGpu(app, layout.gpu, Cfg.UpdateInterval)
		go ui.UpdateGpuTemp(app, layout.gpuTemp, Cfg.UpdateInterval)
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

func updateConfigVars(k *koanf.Koanf, f *file.File) {
	slog.Debug("Loading variables from config file at `" + CONFIG_FILENAME + "` ...")
	if err := k.Load(f, toml.Parser()); err != nil {
		slog.Error("Could not load config file! " + err.Error())
	}
	Cfg = &ConfigVars{
		Debug:          k.Bool("Debug"),
		UpdateInterval: k.Duration("UpdateInterval") * time.Millisecond,
		Celsius:        k.Bool("Celsius"),
		GroupProcesses: k.Bool("GroupProcesses"),
		// TODO: detect AMD / nvidia gpus automatically and override??
		EnableGPU: k.Bool("EnableGPU"),
		EnableUI:  k.Bool("EnableUI"),
	}
	slog.Info("Successfully updated configuration variables from `" + CONFIG_FILENAME + "`")
}

func writeConfigFile() {
	slog.Debug("Creating new config file named `" + CONFIG_FILENAME + "` ...")
	f, err := os.Create(CONFIG_FILENAME)
	if err != nil {
		slog.Error("Could not create config file! " + err.Error())
	}

	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			slog.Error("Could not close config file! " + err.Error())
		}
	}(f)

	updateIntervalStr, _, _ := strings.Cut(Cfg.UpdateInterval.String(), "ms")
	fileData := []byte("Debug=" + strconv.FormatBool(Cfg.Debug) + "\n\n" +
		"# Set how frequently to update the UI (in milliseconds - 1000ms equals 1 second)\n" +
		"UpdateInterval=" + updateIntervalStr + "\n\n" +
		"# Temperature units - `true` for Celsius, `false` for Fahrenheit\n" +
		"Celsius=" + strconv.FormatBool(Cfg.Celsius) + "\n\n" +
		"# Enable or disable grouping of processes in Processes table (true or false)\n" +
		"GroupProcesses=" + strconv.FormatBool(Cfg.GroupProcesses) + "\n\n" +
		"# Enable or disable GPU activity and temperature boxes (true or false)\n" +
		"EnableGPU=" + strconv.FormatBool(Cfg.EnableGPU) + "\n\n" +
		"# this is for debugging... set to false if you want to read startup / setup logs\n" +
		"EnableUI=" + strconv.FormatBool(Cfg.EnableUI) + "\n\n")
	if _, err := f.Write(fileData); err != nil {
		slog.Error("Failed to write config file! " + err.Error())
	}

	slog.Info("Successfully wrote config file to " + CONFIG_FILENAME)
}

func main() {
	app = tview.NewApplication()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(logHandler))

	k := koanf.New(".")
	f := file.Provider(CONFIG_FILENAME)

	// If the config file exists, update `Cfg` using updateConfigVars()
	if _, err := os.Stat(CONFIG_FILENAME); err == nil {
		// load config values from file and start the app
		updateConfigVars(k, f)

		if Cfg.EnableUI {
			// If the text UI is enabled, run the app. Otherwise, don't start it.
			//	This is mostly for debugging. Eventually I'll log to file... but not today
			startApp(app)
		} else {
			slog.Info("Did not start app - EnableUI is " +
				strconv.FormatBool(Cfg.EnableUI) + " !")
		}
	} else {
		// if a config file does not exist, use default values and start the app anyway
		slog.Info("`" + CONFIG_FILENAME + "` does not exist!")
		writeConfigFile()
		time.Sleep(time.Second * 4)

		if Cfg.EnableUI {
			startApp(app)
		} else {
			slog.Info("Did not start app - EnableUI is " + strconv.FormatBool(Cfg.EnableUI) + " !")
		}
	}
}
